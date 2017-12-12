const grpc = require('grpc');
const path = require('path');
const fs = require('fs');
const mkdirp = require('mkdirp');

const protoFile = path.join(__dirname, '../proto/compress.proto');
const jsFile = path.join(__dirname, '../assets/lodash.min.js');
const imgFile = path.join(__dirname, '../assets/banner.png');
const distPath = path.join(__dirname, '../assets/dist');

const compress = grpc.load(protoFile).compress;
const client = new compress.Compress('127.0.0.1:3016', grpc.credentials.createInsecure());

mkdirp.sync(distPath);

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

function doBrotli() {
  const buf = fs.readFileSync(jsFile);
  const request = new compress.CompressRequest();
  request.setType(compress.Type.BROTLI);
  request.setData(new Uint8Array(buf));
  doRequest(request).then((res) => {
    fs.writeFileSync(`${distPath}/lodash.br`, res.data);
  }).catch(console.error);
}

function doGzip() {
  const buf = fs.readFileSync(jsFile);
  const request = new compress.CompressRequest();
  request.setType(compress.Type.GZIP);
  request.setData(new Uint8Array(buf));
  doRequest(request).then((res) => {
    fs.writeFileSync(`${distPath}/lodash.zip`, res.data);
  }).catch(console.error);
}

function doWebp() {
  const buf = fs.readFileSync(imgFile);
  const request = new compress.CompressRequest();
  request.setType(compress.Type.WEBP);
  request.setData(new Uint8Array(buf));
  request.setQuality(75);
  doRequest(request).then((res) => {
    fs.writeFileSync(`${distPath}/banner.webp`, res.data);
  }).catch(console.error);
}

function doJepg() {
  const buf = fs.readFileSync(imgFile);
  const request = new compress.CompressRequest();
  request.setType(compress.Type.JPEG);
  request.setData(new Uint8Array(buf));
  request.setQuality(75);
  doRequest(request).then((res) => {
    fs.writeFileSync(`${distPath}/banner.jpeg`, res.data);
  }).catch(console.error);
}

function doPNG() {
  const buf = fs.readFileSync(imgFile);
  const request = new compress.CompressRequest();
  request.setType(compress.Type.PNG);
  request.setData(new Uint8Array(buf));
  doRequest(request).then((res) => {
    fs.writeFileSync(`${distPath}/banner.png`, res.data);
  }).catch(console.error); 
}

doBrotli();
doGzip();
doWebp();
doJepg();
doPNG();
