# tiny 

此程序主要提供图片的转换处理以及文本的压缩，有`HTTP`与`GRPC`的调用方式，建议配合[tiny-site](https://github.com/vicanso/tiny-site)使用。

- 图片支持`webp` `jpeg` `png`
- 数据压缩支持`brotli` `gzip`两种

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

## quality

- gzip: `-1` - `9`
- brotli: `0` - `11`
- guetzli: `84` - 
- jpeg: `0` - ?
- webp: `0` - ?
- png: `0` - ?

## docker

### build

```bash
docker build -t vicanso/tiny .
```

### run

```bash
docker run -d --restart=always -p 3015:3015 -p 3016:3016 vicanso/tiny
```

## example

- `query.url` 需要做压缩的源数据地址
- `query.output` 输出类型，可以选择png, webp, jpeg
- `query.width` 图片转换后的宽度，如果不设置，自适应
- `query.height` 图片转换后的高度，如果不设置，自适应
- `query.quality` 图片压缩处理时的质量，对于`webp`，`0`表示无损。对于`brotli`，如果为`0`表示默认值`9`。对于`gzip`，如果为`0`表示使用默认压缩级别。


```bash
curl 'http://127.0.0.1:3015/optim?output=webp&url=http://oidmt881u.bkt.clouddn.com/mac.jpg&quality=30'
```