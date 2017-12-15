package shadow

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	pb "../proto"
	"github.com/buger/jsonparser"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// 获取query参数
func getQuery(query map[string][]string, key string) string {
	data := query[key]
	if data == nil {
		return ""
	}
	return strings.Join(data, "")
}

// 读取数据，根据请求的url或者base64数据
func getData(req *http.Request) ([]byte, string, error) {
	query := req.URL.Query()
	url := getQuery(query, "url")
	contentType := ""
	if len(url) != 0 {
		c := &http.Client{
			Timeout: 10 * time.Second,
		}
		res, err := c.Get(url)
		if err != nil {
			return nil, contentType, err
		}
		contentType = res.Header.Get("Content-Type")
		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)
		return data, contentType, err
	}
	body, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil, contentType, err
	}
	base64Data, err := jsonparser.GetString(body, "base64")
	if err != nil {
		return nil, contentType, err
	}
	data, err := base64.StdEncoding.DecodeString(base64Data)
	return data, contentType, err
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

func pingServe(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

// 图片压缩处理（保持原有尺寸，调整质量）
func optimServe(w http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL)
	w.Header().Set("Cache-Control", "no-cache")
	query := req.URL.Query()

	alg := getQuery(query, "type")
	imageType := getQuery(query, "imageType")

	tmpWidth, _ := strconv.Atoi(getQuery(query, "width"))
	width := uint32(tmpWidth)
	tmpHeight, _ := strconv.Atoi(getQuery(query, "height"))
	height := uint32(tmpHeight)
	tmpQuality, _ := strconv.Atoi(getQuery(query, "quality"))
	quality := uint32(tmpQuality)

	data, contentType, err := getData(req)
	header := w.Header()
	if err != nil {
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "load data faial"}`))
		return
	}
	var newBuf []byte
	switch alg {
	default:
		newBuf, err = doGzip(data)
		header.Set("Content-Encoding", "gzip")
	case "brotli":
		newBuf, err = doBrotli(data, quality)
		header.Set("Content-Encoding", "br")
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
		header.Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message": "compress data fail"}`))
		return
	}

	header.Set("Content-Length", strconv.Itoa(len(newBuf)))

	if len(contentType) != 0 {
		header.Set("Content-Type", contentType)
	}
	w.Write(newBuf)
}

// RunHTTP 启动HTTP服务
func RunHTTP() {
	http.HandleFunc("/ping", pingServe)
	http.HandleFunc("/@tiny/optim", optimServe)
}

// RunGRPC 启动GRPC服务
func RunGRPC() *grpc.Server {
	s := grpc.NewServer()
	pb.RegisterCompressServer(s, &server{})
	reflection.Register(s)
	return s
}
