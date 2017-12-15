# compression

## 使用建议

对于图片，我的建议是尽量使用`webp`格式，在客户端判断是否支持`webp`而指定加载的图片格式，而后端实时将图片做转换输出。由于图片都是可以做缓存的，转换后的数据会缓存到`varnish`中，因此实时转换的性能并没有太大的影响，当然也可以在系统发布的时候，直接先将图片转换保存多一份。

在考虑系统是否需要增加`brotli`压缩时，先收集当前用户支持`brotli`的占比，我的系统大概有`50%`的用户是支持的，因此增加`brotli`的支持能减少部分带宽的占用（主要就是省钱，如果是缓存的数据可以选择更高的压缩比），下面是我们现行的处理方式：

- 前置HTTP缓存服务器(varnish)，根据客户端`Accept-Encoding`划分为`brotli`与`gzip`两种
- 后端根据该请求是否可缓存判断使用何种压缩算法，如果该请求可缓存（如新闻列表），则根据`Accept-Encoding`选择压缩算法（因为在varnish缓存的是压缩数据，因此压缩一次之后，在后续缓存期内，可多次使用）。而不可缓存的请求，统一使用`gzip`以提高性能
- 对于静态文件(css, js等)，在发布时自动生成两种压缩文件，`nginx`根据`Accept-Encoding`自动选择返回对应的静态文件

## poroto gen

### gen go 

```bash
protoc -I proto/ proto/compress.proto --go_out=plugins=grpc:proto
```

## docker

export LD_LIBRARY_PATH=/tiny/lib
指定动态库目录

### build

```bash
# 生成动态库与执行文件
docker run -it --rm -v ~/github/tiny:/tiny golang bash
```

```bash
docker build -t vicanso/tiny .
```

## example

- `query.url` 需要做压缩的源数据地址
- `body.base64` 如果`query.url`为空，则从POST数据中的base64字段中取数据
- `query.type` 压缩类型，对于文本数据，可以为：`brotli`或者`gzip`。图片数据则为：`webp`、`jpeg`或者`png`
- `query.imageType` 对于图片数据，建议指定图片类型，否则有可能图片解码失败
- `query.width` 图片转换后的宽度，如果不设置，自适应
- `query.height` 图片转换后的高度，如果不设置，自适应
- `query.quality` 图片压缩处理时的质量，对于`webp`，`0`表示无损。对于`brotli`，如果为`0`表示默认值`9`。

```bash
curl 'http://aslant.site/@tiny/optim?url=https://raw.githubusercontent.com/lodash/lodash/4.17.4/dist/lodash.min.js&type=brotli' 

curl 'http://aslant.site/@tiny/optim?url=https://www.joinpay.com/statics/themes/default/images/bigBanner1.png&type=webp&quality=75'
```