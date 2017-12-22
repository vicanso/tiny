package shadow

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	pb "../proto"
	"github.com/buger/jsonparser"
	"github.com/mozillazg/request"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 获取query参数
func getQuery(query *fasthttp.Args, key string) string {
	data := query.Peek(key)
	if data == nil {
		return ""
	}
	return string(data[:])
}

// 读取数据，根据请求的url或者base64数据
func getData(ctx *fasthttp.RequestCtx) ([]byte, string, error) {
	query := ctx.QueryArgs()
	url := getQuery(query, "url")
	contentType := ""
	if len(url) != 0 {
		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		req := request.NewRequest(c)
		req.Headers = map[string]string{
			"Accept-Encoding": "gzip",
		}
		res, err := req.Get(url)
		if err != nil {
			return nil, "", err
		}
		contentType = res.Header.Get("Content-Type")
		contentEncoding := res.Header.Get("Content-Encoding")
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		if contentEncoding == "gzip" {
			r, err := gzip.NewReader(bytes.NewBuffer(data))
			if err != nil {
				return nil, "", err
			}
			data, err = ioutil.ReadAll(r)
			if err != nil {
				return nil, "", err
			}
		}
		return data, contentType, err
	}
	body := ctx.PostBody()
	data, _ := jsonparser.GetString(body, "data")
	var buf []byte
	if len(data) == 0 {
		base64Data, err := jsonparser.GetString(body, "base64")
		if err != nil {
			return nil, contentType, err
		}
		buf, err = base64.StdEncoding.DecodeString(base64Data)
		if err != nil {
			return nil, contentType, err
		}
	} else {
		buf = []byte(data)
	}
	return buf, contentType, nil
}

func grpcCompress(in *pb.CompressRequest) ([]byte, error) {
	var newBuf []byte
	var err error
	alg := in.Type
	buf := in.Data
	switch alg {
	default:
		newBuf, err = doGzip(buf)
	case pb.Type_BROTLI:
		newBuf, err = doBrotli(buf, in.Quality)
	case pb.Type_WEBP:
		newBuf, err = doWebp(buf, in.Width, in.Height, in.Quality, in.ImageType)
	case pb.Type_JPEG:
		newBuf, err = doJPEG(buf, in.Width, in.Height, in.Quality, in.ImageType)
	case pb.Type_PNG:
		newBuf, err = doPNG(buf, in.Width, in.Height, in.ImageType)
	}
	if err != nil {
		return nil, err
	}
	return newBuf, nil
}

// server is used to implement compress.CompressServer.
type server struct{}

func (s *server) Do(ctx context.Context, in *pb.CompressRequest) (*pb.CompressReply, error) {
	buf, err := grpcCompress(in)
	if err != nil {
		return nil, err
	}
	return &pb.CompressReply{
		Data: buf,
	}, nil
}

func pingServe(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("pong")
	ctx.SetConnectionClose()
}

// 图片压缩处理（保持原有尺寸，调整质量）
func optimServe(ctx *fasthttp.RequestCtx) {
	log.Printf("%s %s %s", ctx.RemoteAddr(), ctx.Method(), ctx.URI())
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	query := ctx.QueryArgs()

	alg := getQuery(query, "type")
	imageType := getQuery(query, "imageType")

	tmpWidth, _ := strconv.Atoi(getQuery(query, "width"))
	width := uint32(tmpWidth)
	tmpHeight, _ := strconv.Atoi(getQuery(query, "height"))
	height := uint32(tmpHeight)
	tmpQuality, _ := strconv.Atoi(getQuery(query, "quality"))
	quality := uint32(tmpQuality)

	data, contentType, err := getData(ctx)
	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Error(`{"message": "load data faial"}`, fasthttp.StatusInternalServerError)
	}
	var newBuf []byte
	switch alg {
	default:
		newBuf, err = doGzip(data)
		ctx.Response.Header.Set("Content-Encoding", "gzip")
	case "brotli":
		newBuf, err = doBrotli(data, quality)
		ctx.Response.Header.Set("Content-Encoding", "br")
	case "webp":
		contentType = "image/webp"
		newBuf, err = doWebp(data, width, height, quality, imageType)
	case "jpeg":
		contentType = "image/jpeg"
		newBuf, err = doJPEG(data, width, height, quality, imageType)
	case "png":
		contentType = "image/png"
		newBuf, err = doPNG(data, width, height, imageType)
	}

	if err != nil {
		ctx.Response.Header.Set("Content-Type", "application/json")
		ctx.Error(`{"message": "compress data fail"}`, fasthttp.StatusInternalServerError)
		return
	}

	ctx.Response.Header.Set("Content-Length", strconv.Itoa(len(newBuf)))

	if len(contentType) != 0 {
		ctx.Response.Header.Set("Content-Type", contentType)
	}
	ctx.SetBody(newBuf)
}

// HTTPHandler 启动HTTP服务
func HTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/ping":
		pingServe(ctx)
	case "/@tiny/optim":
		optimServe(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}

// RunGRPC 启动GRPC服务
func RunGRPC() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &server{})
	reflection.Register(s)
	return s
}
