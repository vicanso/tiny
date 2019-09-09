# tiny

提供图片的转换处理以及文本的压缩，有`HTTP`与`GRPC`的调用方式，建议配合[tiny-site](https://github.com/vicanso/tiny-site)使用。

- 图片支持`webp`, `jpeg`, `png`
- 数据压缩支持`brotli`, `gzip`, `snappy`, `lz4`, `zstd`

## 编译proto

```bash
make protoc
```

## 启动

```bash
docker run -d --restart=always \
  -p 7001:7001 \
  -p 7001:7002 \
  --name=tiny \
  vicanso/tiny:elton
```

### 示例

以brotli方式压缩文件：

```bash
curl 'http://127.0.0.1:7001/texts/optim?output=br&quality=11&url=https://cdn.staticfile.org/jquery/3.4.1/jquery.min.js'
```

以POST的形式指定文本压缩：

```bash
curl -XPOST -H 'Content-Type:application/json' -d '{
	"data": "strings.........",
	"output": "gzip"
}' 'http://127.0.0.1:7001/texts/optim'
```

将png转换为webp:

```bash
curl 'http://127.0.0.1:7001/images/optim?output=webp&quality=80&url=https://www.baidu.com/img/bd_logo1.png'
```

以POST的形式指定图片转换：
```bash
curl -XPOST -H 'Content-Type:application/json' -d '{
	"data": "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",
	"source": "jpeg",
	"output": "jpeg",
	"quality": 80,
	"width": 60,
	"height": 0
}' 'http://127.0.0.1:7001/images/optim'
```

## 客户端

