# tiny

提供图片的转换处理以及文本的压缩，有`HTTP`与`GRPC`的调用方式，建议配合[tiny-site](https://github.com/vicanso/tiny-site)使用。

- `png` PNG的优化处理使用[pngquant](https://github.com/kornelski/pngquant)
- `jpeg` JEPG的优化处理使用[mozjpeg](https://github.com/mozilla/mozjpeg)
- `avif` AVIF的优化处理使用[cavif](https://github.com/kornelski/cavif-rs)

- 图片输出支持`webp`, `jpeg`, `png`, `avif`
- 数据压缩输出支持`brotli`, `gzip`, `snappy`, `lz4`, `zstd`

## 编译proto

需要先安装`protoc-gen-gofast`：

```bash
go get -d github.com/gogo/protobuf/protoc-gen-gofast
```

```bash
make protoc
```

## 启动

```bash
docker run -d --restart=always \
  -p 7001:7001 \
  -p 7002:7002 \
  --name=tiny \
  vicanso/tiny
```

其中7001提供HTTP服务，7002提供GRPC服务，默认http body限制为1MB，如果需要调整，可通过ENV来调整，如`TINY_BODY_PARSER_LIMIT=10MB`

### 示例

以brotli方式压缩文件（需要注意，只有HTTPS或者以IP形式打开，chrome才支持br）：

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

tiny客户端主要用于图片预处理，可以将指定目录下的所有图片压缩优化，参数如下：

```bash
Usage of tiny:
  -filter string
    	filter regexp for image (default ".(png|jpg|jpeg)$")
  -jpeg int
    	the quality of jpeg, it should be >= 0 and <= 100 (default 80)
  -png int
    	the quality of png, it should be >= 0 and <= 100 (default 90)
  -server string
    	grpc server address (default "tiny.aslant.site:7002")
  -source string
    	search path (default ".")
  -target string
    	optim target path, new image will save to this path
  -webp int
    	the quality of webp, it should be >= 0 and <= 100
```

平时常用中，建议启动自己的tiny server，将`-server`参数指定为该服务，避免网络传输带来的压缩延迟。使用默认参数将`/Downloads`目录下的图片压缩优化并保存至`/tmp/images`中：

```
tiny -source=/Downloads -target=/tmp/images
546 / 546 [----------------------------------------------------------------------] 100.00% 44 p/s
********************************TINY********************************
Optimize images is done, use:12.7260134s
Success(538) Fail(8)
Space size reduce from 28 MB to 7.6 MB
Fails: /Downloads/res_1566890283138/assets/imgs/bg_erweima.png rpc error: code = Unknown desc = message=data can not be nil
/Downloads/res_1566890283138/assets/imgs/breakOrWith.png rpc error: code = Unknown desc = png: invalid format: not a PNG file
/Downloads/res_1566890283138/assets/imgs/default_ishare.png rpc error: code = Unknown desc = png: invalid format: not a PNG file
/Downloads/res_1566890283138/assets/screen/750x1134.png rpc error: code = Unknown desc = message=data can not be nil
/Downloads/res_1566890283138/assets/screen/Default-568h@2x.png rpc error: code = Unknown desc = message=data can not be nil
/Downloads/res_1566890283138/assets/screen/Default-iOS11-812h@3x.png rpc error: code = Unknown desc = message=data can not be nil
/Downloads/res_1566890283138/assets/screen/splash-port-pdpi.png rpc error: code = Unknown desc = message=data can not be nil
/Downloads/res_1566890283138/assets/screen/splash-port-xhdpi.png rpc error: code = Unknown desc = message=data can not be nil
********************************TINY********************************
```
