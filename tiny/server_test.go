package tiny

import (
	"bytes"
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/h2non/gock"
	pb "github.com/vicanso/tiny/proto"
)

const (
	jpegURL = "http://127.0.0.1/ai.jpeg"
	pngURL  = "http://127.0.0.1/ai.png"
	jsURL   = "http://127.0.0.1/debounce.js"
	jsFile  = "../assets/debounce.js"
)

func mockJs() {
	b, _ := ioutil.ReadFile("../assets/debounce.js")
	gzipData, _ := DoGzip(b, 0)
	gock.New(jsURL).
		Reply(200).
		SetHeader("Content-Type", "text/javascript").
		SetHeader("Content-Encoding", "gzip").
		Body(bytes.NewReader(gzipData))
}

func mockJpeg() {
	b, _ := ioutil.ReadFile("../assets/ai.jpeg")
	gock.New(jpegURL).
		Reply(200).
		SetHeader("Content-Type", "image/jpeg").
		Body(bytes.NewReader(b))
}

func mockPng() {
	b, _ := ioutil.ReadFile("../assets/ai.png")
	gock.New(pngURL).
		Reply(200).
		SetHeader("Content-Type", "image/png").
		Body(bytes.NewReader(b))
}

func TestGetOptimParams(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/optim?output=png&url=a&width=1&height=2&quality=3", nil)
	params, err := getOptimParams(r)
	if err != nil {
		t.Fatalf("get params fail, %v", err)
	}
	if params.Output != "png" ||
		params.URL != "a" ||
		params.Width != 1 ||
		params.Height != 2 ||
		params.Quality != 3 {
		t.Fatalf("get params fail")
	}
}

func TestGetOptimData(t *testing.T) {
	defer gock.Off()
	mockJpeg()
	body, err := getOptimData(jpegURL)
	if err != nil {
		t.Fatalf("get optim data fail, %v", err)
	}
	if body.Type != "image/jpeg" || len(body.Data) != 6641 {
		t.Fatalf("get optim image data fail")
	}

	mockJs()
	body, err = getOptimData(jsURL)
	if err != nil {
		t.Fatalf("get optim data fail, %v", err)
	}
	if body.Type != "text/javascript" || len(body.Data) != 6661 {
		t.Fatalf("get optim text data fail")
	}
}

func TestResErr(t *testing.T) {
	err := errors.New("test error")
	w := httptest.NewRecorder()
	resRawError(w, err)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("res error status code fail")
	}
	if string(w.Body.Bytes()) != `{"message":"test error"}` {
		t.Fatalf("res error data fail")
	}
}

func TestPing(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "http://127.0.0.1/ping", nil)
	w := httptest.NewRecorder()
	s := HTTPServer{}
	s.Ping(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("ping status code fail")
	}
	if string(w.Body.Bytes()) != "pong" {
		t.Fatalf("ping response fail")
	}
}

func TestOptim(t *testing.T) {
	t.Run("url is nil", func(t *testing.T) {
		url := "http://127.0.0.1/optim?output=webp&quality=50"
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("no url should return error")
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		url := "http://127.0.0.1/optim?output=webp&quality=50&url=abcd"
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("invalid url should return error")
		}
	})

	t.Run("not support output", func(t *testing.T) {
		defer gock.Off()
		mockJpeg()
		url := "http://127.0.0.1/optim?output=a&quality=50&url=" + jpegURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("not support output should return error")
		}

		mockJs()
		url = "http://127.0.0.1/optim?output=a&quality=0&url=" + jsURL
		r = httptest.NewRequest(http.MethodGet, url, nil)
		w = httptest.NewRecorder()
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("not support output should return error")
		}
	})

	t.Run("get optim data fail", func(t *testing.T) {
		defer gock.Off()
		gock.New(jsURL).
			Reply(500).
			SetHeader("Content-Type", "text/javascript").
			SetHeader("Content-Encoding", "gzip").
			BodyString("")
		url := "http://127.0.0.1/optim?output=gzip&quality=0&url=" + jsURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("get optim data fail should return error")
		}
	})

	t.Run("not support content encoding", func(t *testing.T) {
		defer gock.Off()
		gock.New(jsURL).
			Reply(200).
			SetHeader("Content-Type", "text/javascript").
			SetHeader("Content-Encoding", "br").
			BodyString("")

		url := "http://127.0.0.1/optim?output=gzip&quality=0&url=" + jsURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("get optim data not support content encoding fail should return error")
		}
	})

	t.Run("webp", func(t *testing.T) {
		defer gock.Off()
		mockJpeg()
		url := "http://127.0.0.1/optim?output=webp&quality=50&url=" + jpegURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim image to webp fail")
		}
		if w.HeaderMap["Content-Type"][0] != "image/webp" ||
			w.HeaderMap["Content-Length"][0] != "1028" {
			t.Fatalf("optim image to png fail")
		}
	})

	t.Run("to png", func(t *testing.T) {
		defer gock.Off()
		mockPng()
		url := "http://127.0.0.1/optim?output=png&quality=50&url=" + pngURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim image to png fail")
		}
		if w.HeaderMap["Content-Type"][0] != "image/png" ||
			w.HeaderMap["Content-Length"][0] != "1302" {
			t.Fatalf("optim image to png fail")
		}
	})

	t.Run("to jpeg", func(t *testing.T) {
		defer gock.Off()
		mockJpeg()
		url := "http://127.0.0.1/optim?output=jpeg&quality=50&url=" + jpegURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim image to jpeg fail")
		}
		if w.HeaderMap["Content-Type"][0] != "image/jpeg" ||
			w.HeaderMap["Content-Length"][0] != "2074" {
			t.Fatalf("optim image to jpeg fail")
		}
	})

	t.Run("to guetzli", func(t *testing.T) {
		defer gock.Off()
		mockJpeg()
		// guetzli's quality should gte 84
		url := "http://127.0.0.1/optim?output=guetzli&quality=90&url=" + jpegURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim image to guetzli fail")
		}
		if w.HeaderMap["Content-Type"][0] != "image/jpeg" ||
			w.HeaderMap["Content-Length"][0] != "2158" {
			t.Fatalf("optim image to jpeg fail")
		}
	})

	t.Run("resize", func(t *testing.T) {
		defer gock.Off()
		mockJpeg()
		url := "http://127.0.0.1/optim?output=jpeg&quality=50&width=64&url=" + jpegURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim and resize image to jpeg fail")
		}
		if w.HeaderMap["Content-Type"][0] != "image/jpeg" ||
			w.HeaderMap["Content-Length"][0] != "1151" {
			t.Fatalf("optim and resize image to jpeg fail")
		}
	})

	t.Run("gzip", func(t *testing.T) {
		defer gock.Off()
		mockJs()
		url := "http://127.0.0.1/optim?output=gzip&quality=0&url=" + jsURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim gzip fail")
		}
		if w.HeaderMap["Content-Type"][0] != "text/javascript" ||
			w.HeaderMap["Content-Length"][0] != "2265" ||
			w.HeaderMap["Content-Encoding"][0] != "gzip" {
			t.Fatalf("optim gzip fail")
		}
	})

	t.Run("br", func(t *testing.T) {
		defer gock.Off()
		mockJs()
		url := "http://127.0.0.1/optim?output=br&quality=0&url=" + jsURL
		r := httptest.NewRequest(http.MethodGet, url, nil)
		w := httptest.NewRecorder()
		s := HTTPServer{}
		s.Optim(w, r)
		if w.Code != http.StatusOK {
			t.Fatalf("optim br fail")
		}
		if w.HeaderMap["Content-Type"][0] != "text/javascript" ||
			w.HeaderMap["Content-Length"][0] != "2102" ||
			w.HeaderMap["Content-Encoding"][0] != "br" {
			t.Fatalf("optim br fail")
		}
	})
}

func TestGPRCOptim(t *testing.T) {
	imageBuf, _ := ioutil.ReadFile(jpgFile)
	textBuf, _ := ioutil.ReadFile(jsFile)

	gs := GRPCServer{}
	t.Run("webp", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:      pb.Type_WEBP,
			Data:      imageBuf,
			Quality:   90,
			ImageType: JPEG,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 1890 {
			t.Fatalf("grpc optim webp fail, %v", err)
		}
	})

	t.Run("jpeg", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:      pb.Type_JPEG,
			Data:      imageBuf,
			Quality:   90,
			ImageType: JPEG,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 3452 {
			t.Fatalf("grpc optim jpeg fail, %v", err)
		}
	})

	t.Run("png", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:      pb.Type_PNG,
			Data:      imageBuf,
			Quality:   90,
			ImageType: JPEG,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 1630 {
			t.Fatalf("grpc optim png fail, %v", err)
		}
	})

	t.Run("guetzli", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:      pb.Type_GUETZLI,
			Data:      imageBuf,
			Quality:   90,
			ImageType: JPEG,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 2158 {
			t.Fatalf("grpc optim guetzli fail, %v", err)
		}
	})

	t.Run("resize", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:      pb.Type_JPEG,
			Data:      imageBuf,
			Quality:   90,
			ImageType: JPEG,
			Width:     64,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 1736 {
			t.Fatalf("grpc optim jpeg fail, %v", err)
		}
	})

	t.Run("gzip", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:    pb.Type_GZIP,
			Data:    textBuf,
			Quality: 9,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 2264 {
			t.Fatalf("grcp optim gzip fail, %v", err)
		}
	})

	t.Run("br", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:    pb.Type_BROTLI,
			Data:    textBuf,
			Quality: 9,
		}
		buf, err := gs.Optim(in)
		if err != nil || len(buf) != 2102 {
			t.Fatalf("grcp optim br fail, %v", err)
		}
	})

	t.Run("do", func(t *testing.T) {
		in := &pb.CompressRequest{
			Type:    pb.Type_BROTLI,
			Data:    textBuf,
			Quality: 9,
		}
		reply, err := gs.Do(context.Background(), in)
		if err != nil || reply == nil || len(reply.Data) == 0 {
			t.Fatalf("do grpc fail, %v", err)
		}
	})
}
