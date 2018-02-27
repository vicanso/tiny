package shadow

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/mozillazg/request"
	"github.com/valyala/fasthttp"
	pb "github.com/vicanso/tiny/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func grpcCompress(in *pb.CompressRequest) ([]byte, error) {
	var newBuf []byte
	var err error
	alg := in.Type
	buf := in.Data
	imageType := int(in.ImageType)
	quality := int(in.Quality)
	switch alg {
	default:
		newBuf, err = DoGzip(buf, int(in.Quality))
	case pb.Type_BROTLI:
		newBuf, err = DoBrotli(buf, int(in.Quality))
	case pb.Type_WEBP:
		newBuf, err = DoWebp(buf, imageType, quality)
	case pb.Type_JPEG:
		newBuf, err = DoJPEG(buf, imageType, quality)
	case pb.Type_PNG:
		newBuf, err = DoPNG(buf, imageType, quality)
	}
	if err != nil {
		return nil, err
	}
	return newBuf, nil
}

// Server rpc servere处理
type Server struct{}

// Do 数据处理
func (s *Server) Do(ctx context.Context, in *pb.CompressRequest) (*pb.CompressReply, error) {
	buf, err := grpcCompress(in)
	if err != nil {
		return nil, err
	}
	return &pb.CompressReply{
		Data: buf,
	}, nil
}

// GetGRPCServer 获取GRPC Server
func GetGRPCServer() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &Server{})
	reflection.Register(s)
	return s
}

func pingServe(ctx *fasthttp.RequestCtx) {
	ctx.SetBodyString("pong")
	ctx.SetConnectionClose()
}

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
	if len(data) != 0 {
		body, err := base64.StdEncoding.DecodeString(data)
		if err != nil {
			return nil, contentType, err
		}
		buf = body
	} else {
		buf = body
	}
	return buf, contentType, nil
}

// 图片压缩处理（保持原有尺寸，调整质量）
func optimServe(ctx *fasthttp.RequestCtx) {
	log.Printf("%s %s %s", ctx.RemoteAddr(), ctx.Method(), ctx.URI())
	ctx.Response.Header.Set("Cache-Control", "no-cache")
	query := ctx.QueryArgs()

	alg, _ := strconv.Atoi(getQuery(query, "type"))
	imageType, _ := strconv.Atoi(getQuery(query, "imageType"))

	width, _ := strconv.Atoi(getQuery(query, "width"))
	height, _ := strconv.Atoi(getQuery(query, "height"))
	quality, _ := strconv.Atoi(getQuery(query, "quality"))
	data, contentType, err := getData(ctx)
	if err != nil {
		ctx.Error("load data fail, "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	if len(contentType) != 0 {
		switch contentType {
		case "image/jpeg":
			imageType = JPEG
		case "image/png":
			imageType = PNG
		case "image/webp":
			imageType = WEBP
		}
	}
	in := &pb.CompressRequest{
		Type:      pb.Type(alg),
		ImageType: pb.Type(imageType),
		Width:     uint32(width),
		Height:    uint32(height),
		Quality:   uint32(quality),
		Data:      data,
	}
	buf, err := grpcCompress(in)
	if err != nil {
		ctx.Error("optim fail, "+err.Error(), fasthttp.StatusInternalServerError)
		return
	}
	respHeader := &ctx.Response.Header
	switch alg {
	case GZIP:
		respHeader.Set("Content-Encoding", "gzip")
	case BROTLI:
		respHeader.Set("Content-Encoding", "br")
	case WEBP:
		respHeader.SetContentType("image/webp")
	case JPEG:
		respHeader.SetContentType("image/jpeg")
	case PNG:
		respHeader.SetContentType("image/png")
	}
	ctx.SetBody(buf)
}

// HTTPHandler 获取HTTP处理函数
func HTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/ping":
		pingServe(ctx)
	case "/@tiny/optim":
		optimServe(ctx)
	case "/optim":
		optimServe(ctx)
	default:
		ctx.Error("Unsupported path", fasthttp.StatusNotFound)
	}
}
