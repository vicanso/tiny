# compression

## 压缩对比

| 文件名 | 算法 | 原文件大小 | 压缩后大小 | 压缩比 | 耗时(ms)
| font-awesome.css | brotli | 31000 | 6569 | 21% | 22 |
| font-awesome.css | gzip | 31000 | 6968 | 22% | 2 |
| main.js | brotli | 486495 | 109138 | 22% | 128 |
| main.js | gzip | 486495 | 120284 | 24% | 24 |
| react.js | brotli | 6617 | 2683 | 40% | 9 |
| react.js | gzip | 6617 | 2836 | 42% | 1 |
| styles.css | brotli | 115554 | 20510 | 17% | 28 |
| styles.css | gzip | 115554 | 22804 | 19% | 6 |
| vue.js | brotli | 86676 | 30323 | 34% | 32 |
| vue.js | gzip | 86676 | 31819 | 36% | 8 |


| 原始数据 | brotli(字节) | gzip(字节) | brotli(耗时)  | gzip(耗时)
| 726342 | 169223 | 184711 | 219 | 41

从上面的数据可以看出，使用`brotli(Quality:9)`比`gzip`减少了`15488`字节的数据，压缩时间差不多是`5`倍，但总体的压缩时长还算是比较短。

注：如果`brotli`设置`Quality:10`比`gzip`减少了`25455`字节的数据，但是压缩的时间差不多是`50`倍，如果追求更高的压缩率，可以调整`Quality`

## 使用建议

在考虑系统是否需要增加`brotli`压缩时，先收集当前用户支持`brotli`的占比，我们的系统大概有`50%`的用户是支持的，因此增加`brotli`的支持能减少部分带宽的占用（主要就是省钱），下面是我们现行的处理方式：

- 前置HTTP缓存服务器(varnish)，根据客户端`Accept-Encoding`划分为`brotli`与`gzip`两种
- 后端根据该请求是否可缓存判断使用何种压缩算法，如果该请求可缓存（如新闻列表），则根据`Accept-Encoding`选择压缩算法（因为在varnish缓存的是压缩数据，因此压缩一次之后，在后续缓存期内，可多次使用）。而不可缓存的请求，统一使用`gzip`以提高性能
- 对于静态文件(css, js等)，在发布时自动生成两种压缩文件，`nginx`根据`Accept-Encoding`自动选择返回对应的静态文件

## poroto gen

protoc -I compress/ compress/compress.proto --go_out=plugins=grpc:compress