const Koa = require('koa');
const Router = require('koa-router');
const multer = require('koa-multer');
const LRU = require('lru-cache');
const fs = require('fs');
const grpc = require('grpc');
const path = require('path');
const bodyparser = require('koa-bodyparser');
const uuidv4 = require('uuid/v4');
const serve = require('koa-static-serve');

const config = require('./config');

const {
  promisify,
} = require('util');

const readFile = promisify(fs.readFile);
const writeFile = promisify(fs.writeFile);
const unlink = promisify(fs.unlink);
const stat = promisify(fs.stat);
const protoFile = path.join(__dirname, './proto/compress.proto');
const distPath = path.join(__dirname, './dist');

const compress = grpc.load(protoFile).compress;
const client = new compress.Compress(config.grpcServer, grpc.credentials.createInsecure());

const tmpFileCache = new LRU({
  max: 500,
  dispose: (key, value) => {
    unlink(value.path);
  },
  maxAge: 60 * 60 * 1000,
});
const filePath = '/tmp';
const upload = multer({
  dest: filePath,
  limits: {
    fileSize: 1024 * 1024,
  },
});
const router = new Router();
const app = new Koa();

// 获取数据类型
function getType(mode) {
  switch (mode) {
    case 0:
      return compress.Type.GZIP;
    case 1:
      return compress.Type.BROTLI;
    case 2:
      return compress.Type.JPEG;
    case 3:
      return compress.Type.PNG;
    case 4:
      return compress.Type.WEBP;
    case 5:
      return compress.Type.GUETZLI;
    default:
      return compress.Type.GZIP;
  }
}

// 执行grpc调用
function doRequest(request) {
  return new Promise((resolve, reject) => {
    client.do(request, (err, res) => {
      if (err) {
        reject(err);
      } else {
        resolve(res);
      }
    });
  });
}

router.get('/ping', (ctx) => {
  ctx.body = 'pong';
});

// 文件上传
router.post('/api/upload', upload.single('file'), (ctx) => {
  const {
    filename,
    mimetype,
    path,
  } = ctx.req.file;
  tmpFileCache.set(filename, {
    path,
    type: mimetype,
  });
  ctx.body = {
    file: filename,
  };
});

// 文件预览 
router.get('/api/file/:id', async (ctx) => {
  const {
    id,
  } = ctx.params;
  const data = tmpFileCache.get(id);
  if (!data) {
    throw new Error('the image is out of date, please upload again');
  }
  const buf = await readFile(data.path);
  ctx.set('Content-Type', data.type);
  ctx.body = buf;
});

router.get('/api/download/:id', async (ctx) => {
  const {
    id,
  } = ctx.params;
  const data = tmpFileCache.get(id);
  if (!data) {
    throw new Error('the image is out of date, please upload again');
  }
  const buf = await readFile(data.path);
  const filename = `${id}${data.ext}`;
  ctx.set('Content-Disposition', `attachment; filename=${filename}`);
  ctx.body = buf;
});

router.post('/api/tiny', async (ctx) => {
  const {
    mode,
    file,
    quality,
  } = ctx.request.body;
  const data = tmpFileCache.get(file);
  if (!data) {
    throw new Error('the image is out of date, please upload again');
  } 
  const buf = await readFile(data.path);
  const request = new compress.CompressRequest();
  request.setType(getType(mode));
  request.setData(new Uint8Array(buf));
  request.setQuality(quality);

  const res = await doRequest(request);
  const id = uuidv4().replace(/-/g, '');
  const targetData = {
    path: path.join(filePath, id),
  };
  switch (mode) {
    case 0:
      targetData.ext = '.gzip';
      targetData.encoding = 'gzip';
      targetData.type = data.type; 
      break;
    case 1:
      targetData.ext = '.br';
      targetData.encoding = 'br';
      targetData.type = data.type;
      break;
    case 2:
    case 5:
      targetData.ext = '.jpg';
      targetData.type = 'image/jpeg';
      break;
    case 3:
      targetData.ext = '.png';
      targetData.type = 'image/png';
      break;
    case 4:
      targetData.ext = '.webp';
      targetData.type = 'image/webp';
      break;
    default:
      break;
  }
  switch (data.type) {
    case 'image/jpeg':
      request.setImageType(compress.Type.JPEG);
      break;
    case 'image/png':
      request.setImageType(compress.Type.PNG);
      break; 
    case 'image/webp':
      request.setImageType(compress.Type.PNG);
      break;
  }
  const statInfo = await stat(data.path);
  await writeFile(path.join(filePath, id), res.data);
  tmpFileCache.set(id, targetData);
  ctx.body = {
    file: id,
    size: res.data.length,
    originalSize: statInfo.size,
  };
});

if (config.env !== 'development') {
  router.get('/', async (ctx) => {
    const file = path.join(distPath, 'index.html');
    const buf = await readFile(file);
    ctx.set('Content-Type', 'text/html; charset=utf-8');
    ctx.set('Cache-Control', 'public, max-age=60');
    ctx.body = buf;
  });
  app.use(serve(distPath, {
    maxAge: 72 * 60 * 60,
    sMaxAge: 600,
    dotfiles: 'allow',
    denyQuerystring: true,
    etag: false,
    lastModified: false,
    '404': 'next',
    extname: ['.html'],
  }));
}

app.use(async (ctx, next) => {
  try {
    await next();
  } catch (err) {
    ctx.status = 500;
    ctx.body = {
      message: err.message,
    };
  }
});

app
  .use(bodyparser())
  .use(router.routes())
  .use(router.allowedMethods());

app.listen(5018);
