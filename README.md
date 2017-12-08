# compression

## 使用建议

在考虑系统是否需要增加`brotli`压缩时，先收集当前用户支持`brotli`的占比，我们的系统大概有`50%`的用户是支持的，因此增加`brotli`的支持能减少部分带宽的占用（主要就是省钱），下面是我们现行的处理方式：

- 前置HTTP缓存服务器(varnish)，根据客户端`Accept-Encoding`划分为`brotli`与`gzip`两种
- 后端根据该请求是否可缓存判断使用何种压缩算法，如果该请求可缓存（如新闻列表），则根据`Accept-Encoding`选择压缩算法（因为在varnish缓存的是压缩数据，因此压缩一次之后，在后续缓存期内，可多次使用）。而不可缓存的请求，统一使用`gzip`以提高性能
- 对于静态文件(css, js等)，在发布时自动生成两种压缩文件，`nginx`根据`Accept-Encoding`自动选择返回对应的静态文件

## poroto gen

### gen go 

```bash
protoc -I compress/ compress/compress.proto --go_out=plugins=grpc:compress
```

### gen nodejs

```bash
npm install -g grpc-tools

grpc_tools_node_protoc --js_out=import_style=commonjs,binary:node --grpc_out=node --plugin=protoc-gen-grpc=`which grpc_tools_node_protoc_plugin` compress/compress.proto
```