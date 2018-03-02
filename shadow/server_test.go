package shadow

import (
	"context"
	"encoding/base64"
	"io/ioutil"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
	pb "github.com/vicanso/tiny/proto"
)

func doGrpc(file string, imageType, dstType pb.Type, quality int) (*pb.CompressReply, error) {
	server := &Server{}
	buf, _ := ioutil.ReadFile(file)
	in := &pb.CompressRequest{
		Type:    dstType,
		Data:    buf,
		Quality: uint32(quality),
	}
	d := time.Now().Add(3 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	return server.Do(ctx, in)
}

func TestDoBrotli(t *testing.T) {
	out, err := doGrpc("../assets/angular.min.js", 0, pb.Type_BROTLI, 9)
	if err != nil {
		t.Fatalf("grpc do fail, %v", err)
	}
	size := len(out.Data)
	if size != brotliSize {
		t.Fatalf("do brotli fail, size:%d", size)
	}
}

func TestDoGzip(t *testing.T) {
	out, err := doGrpc("../assets/angular.min.js", 0, pb.Type_GZIP, 0)
	if err != nil {
		t.Fatalf("grpc do fail, %v", err)
	}
	size := len(out.Data)
	if size != gzipSize {
		t.Fatalf("do gzip fail, size:%d", size)
	}
}

func TestDoWebp(t *testing.T) {
	out, err := doGrpc("../assets/mac.jpg", pb.Type_JPEG, pb.Type_WEBP, 75)
	if err != nil {
		t.Fatalf("grpc do fail, %v", err)
	}
	size := len(out.Data)
	if size != webpSize {
		t.Fatalf("do webp fail, size:%d", size)
	}
}

func TestDoJpeg(t *testing.T) {
	out, err := doGrpc("../assets/mac.jpg", pb.Type_JPEG, pb.Type_JPEG, 90)
	if err != nil {
		t.Fatalf("grpc do fail, %v", err)
	}
	size := len(out.Data)
	if size != jpegSize {
		t.Fatalf("do jpeg fail, size:%d", size)
	}
}

func TestDoPng(t *testing.T) {
	out, err := doGrpc("../assets/mac.jpg", pb.Type_JPEG, pb.Type_PNG, 80)
	if err != nil {
		t.Fatalf("grpc do fail, %v", err)
	}
	size := len(out.Data)
	if size != pngSize {
		t.Fatalf("do png fail, size:%d", size)
	}
}

func getCtx(file, query string) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetMethod("POST")
	ctx.Request.SetRequestURI("/optim?" + query)
	buf, _ := ioutil.ReadFile(file)
	str := base64.StdEncoding.EncodeToString(buf)
	body := []byte(`{"data": "` + str + `"}`)
	ctx.Request.SetBody(body)
	return ctx
}

func TestHttpDoCompress(t *testing.T) {
	ctx := getCtx("../assets/angular.min.js", "type=1&quality=9")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != brotliSize {
		t.Fatalf("brotli compress by http fail")
	}

	ctx = getCtx("../assets/angular.min.js", "type=0&quality=0")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != gzipSize {
		t.Fatalf("gzip compress by http fail")
	}
}

func TestHttpDoImageConvert(t *testing.T) {
	ctx := getCtx("../assets/mac.jpg", "type=2&imageType=2&quality=90")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != jpegSize {
		t.Fatalf("convert jpg to jpg(90) by http fail")
	}

	ctx = getCtx("../assets/mac.jpg", "type=3&imageType=2&quality=80")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != pngSize {
		t.Fatalf("convert jpg to png by http fail")
	}

	ctx = getCtx("../assets/mac.jpg", "type=4&imageType=2&quality=75")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != webpSize {
		t.Fatalf("convert jpg to webp by http fail")
	}

	ctx = &fasthttp.RequestCtx{}
	ctx.Request.SetRequestURI("/optim?type=2&imageType=2&quality=90&url=http%3A%2F%2Foidmt881u.bkt.clouddn.com%2Fmac.jpg")
	HTTPHandler(ctx)
	if len(ctx.Response.Body()) != jpegSize {
		t.Fatalf("convert jpg to jpg(90) by http fail")
	}
}
