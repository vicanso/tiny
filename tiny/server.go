package tiny

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"log"

	pb "github.com/vicanso/tiny/proto"
)

const (
	// ContentType content type
	ContentType = "Content-Type"
	// ContentEncoding content encoding
	ContentEncoding = "Content-Encoding"
	// ContentLength content length
	ContentLength = "Content-Length"
	// CacheControl cache control
	CacheControl = "Cache-Control"
	// ImageWeb image webp
	ImageWeb = "image/webp"
	// ImagePng image png
	ImagePng = "image/png"
	// ImageJpeg image jpeg
	ImageJpeg = "image/jpeg"
)

type (
	// HTTPServer http server
	HTTPServer struct{}
	// HTTPError http error
	HTTPError struct {
		Code    int
		Message string
	}
	// OptimParms optim params
	OptimParms struct {
		Output  string `json:"output,omitempty"`
		URL     string `json:"url,omitempty"`
		Clip    string `json:"clip,omitempty"`
		Width   int    `json:"width,omitempty"`
		Height  int    `json:"height,omitempty"`
		Quality int    `json:"quality,omitempty"`
	}
	// OptimBody optim body
	OptimBody struct {
		Data []byte
		Type string
	}
	// GRPCServer grpc server
	GRPCServer struct{}
)

func getOptimParams(r *http.Request) (params *OptimParms, err error) {
	m := make(map[string]interface{})
	convertKeys := "width height quality"
	for k, v := range r.URL.Query() {
		if strings.Contains(convertKeys, k) {
			m[k], _ = strconv.Atoi(v[0])
		} else {
			m[k] = v[0]
		}
	}
	buf, err := json.Marshal(m)
	if err != nil {
		return
	}
	params = &OptimParms{}
	err = json.Unmarshal(buf, params)
	return
}

func getOptimData(url string) (body *OptimBody, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("get optim data fail")
		return
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	header := resp.Header
	encoding := header.Get("Content-Encoding")
	if encoding != "" {
		if encoding != "gzip" {
			err = errors.New("not support the encoding")
			return
		}
		buf, err = DoGunzip(buf)
		if err != nil {
			return
		}
	}
	body = &OptimBody{
		Data: buf,
		Type: header.Get(ContentType),
	}
	return
}

func resErr(w http.ResponseWriter, err *HTTPError) {
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set(ContentType, "application/json; charset=utf-8")
	if err.Code != 0 {
		w.WriteHeader(err.Code)
	}
	m := map[string]string{
		"message": err.Message,
	}
	buf, _ := json.Marshal(m)
	w.Write(buf)
}

func resRawError(w http.ResponseWriter, err error) {
	resErr(w, &HTTPError{
		Code:    http.StatusBadRequest,
		Message: err.Error(),
	})
}

// Ping the health check function
func (s *HTTPServer) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

// Optim optim the image
func (s *HTTPServer) Optim(w http.ResponseWriter, r *http.Request) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		log.Println("optim fail, " + err.Error())
	}()
	var err error
	params, err := getOptimParams(r)
	if err != nil {
		resRawError(w, err)
		return
	}
	if params.Output == "" || params.URL == "" {
		resErr(w, &HTTPError{
			Code:    http.StatusBadRequest,
			Message: "output and url can not be nil",
		})
		return
	}
	body, err := getOptimData(params.URL)
	if err != nil {
		resRawError(w, err)
		return
	}
	header := w.Header()
	output := 0
	switch params.Output {
	case "webp":
		output = WEBP
		header.Set(ContentType, ImageWeb)
	case "png":
		output = PNG
		header.Set(ContentType, ImagePng)
	case "jpeg":
		fallthrough
	case "jpg":
		output = JPEG
		header.Set(ContentType, ImageJpeg)
	case "guetzli":
		output = GUETZLI
		header.Set(ContentType, ImageJpeg)
	case "gzip":
		output = GZIP
		header.Set(ContentEncoding, "gzip")
	case "br":
		output = BROTLI
		header.Set(ContentEncoding, "br")
	default:
		// not support output is return error
		resRawError(w, errors.New("not support the output type"))
		return
	}
	imageType := 0
	switch body.Type {
	case "image/png":
		imageType = PNG
	case "image/jpeg":
		imageType = JPEG
	}

	var buf []byte
	quality := params.Quality
	// optim image
	if imageType > 0 {
		if params.Clip != "" {
			clipType := ClipCenter
			// 暂时仅支持两种方式
			if params.Clip == "leftTop" {
				clipType = ClipLeftTop
			}
			buf, err = DoClip(body.Data, clipType, imageType, quality, params.Width, params.Height, output)
		} else if params.Height != 0 || params.Width != 0 {
			buf, err = DoResize(body.Data, imageType, quality, params.Width, params.Height, output)
		} else {
			switch output {
			case WEBP:
				buf, err = DoWebp(body.Data, imageType, quality)
			case JPEG:
				buf, err = DoJPEG(body.Data, imageType, quality)
			case PNG:
				buf, err = DoPNG(body.Data, imageType, quality)
			case GUETZLI:
				buf, err = DoGUEZLI(body.Data, imageType, quality)
			}
		}
	} else {
		// optim text (gzip br)
		header.Set(ContentType, body.Type)
		switch output {
		case GZIP:
			buf, err = DoGzip(body.Data, quality)
		case BROTLI:
			buf, err = DoBrotli(body.Data, quality)
		}
	}
	if err != nil {
		resRawError(w, err)
		return
	}
	header.Set(ContentLength, strconv.Itoa(len(buf)))
	header.Set(CacheControl, "public, max-age=86400")
	_, err = w.Write(buf)
	if err != nil {
		fmt.Sprintln("write buffer fail,", err)
	}
}

// Optim optim by grpc
func (gs *GRPCServer) Optim(in *pb.CompressRequest) ([]byte, error) {
	var newBuf []byte
	var err error
	alg := in.Type
	buf := in.Data
	imageType := int(in.ImageType)
	quality := int(in.Quality)
	width := int(in.Width)
	height := int(in.Height)
	clipType := int(in.ClipType)
	switch alg {
	default:
		newBuf, err = DoGzip(buf, quality)
	case pb.Type_BROTLI:
		newBuf, err = DoBrotli(buf, quality)
	case pb.Type_WEBP:
		fallthrough
	case pb.Type_JPEG:
		fallthrough
	case pb.Type_PNG:
		fallthrough
	case pb.Type_GUETZLI:
		if clipType != ClipNone {
			newBuf, err = DoClip(buf, clipType, imageType, quality, width, height, int(alg))
		} else {
			newBuf, err = DoResize(buf, imageType, quality, width, height, int(alg))
		}
	}

	if err != nil {
		return nil, err
	}
	return newBuf, nil
}

// Do grpc server do
func (gs *GRPCServer) Do(ctx context.Context, in *pb.CompressRequest) (*pb.CompressReply, error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		err, ok := r.(error)
		if !ok {
			err = fmt.Errorf("%v", r)
		}
		log.Println("optim fail, " + err.Error())
	}()
	buf, err := gs.Optim(in)
	if err != nil {
		return nil, err
	}
	return &pb.CompressReply{
		Data: buf,
	}, nil
}
