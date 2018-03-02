package shadow

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

const (
	// brotliSIZE br压缩后的数据大小
	brotliSize = 55893
	// gzipSize gzip压缩后的数据大小
	gzipSize       = 59656
	webpSize       = 124770
	jpegSize       = 222640
	pngSize        = 415628
	webp480Size    = 31392
	jpeg480Size    = 57797
	png480Size     = 751897
	webpResizeSize = 207480
	jpegGuezliSize = 160977
)

func compressImage(t *testing.T, fn func([]byte, int, int) ([]byte, error), imageType, quality int, file, alg, dst string) []byte {
	buf, _ := ioutil.ReadFile(file)
	newBuf, err := fn(buf, imageType, quality)
	if err != nil {
		t.Fatalf(alg+" fail, %v", err)
		return nil
	}
	log.Printf(alg+" compress, original:%d compress:%d", len(buf), len(newBuf))
	ioutil.WriteFile(dst, newBuf, os.ModePerm)
	return newBuf
}

func TestWebp(t *testing.T) {
	buf := compressImage(t, DoWebp, PNG, 0, "../assets/fluidicon.png", "webp", "../assets/compress/fluidicon.webp")
	if len(buf) != 21302 {
		t.Fatalf("convert png to webp fail")
	}
	buf = compressImage(t, DoWebp, JPEG, 75, "../assets/mac.jpg", "webp", "../assets/compress/mac.webp")
	if len(buf) != webpSize {
		t.Fatalf("convert jpeg to webp fail")
	}
	buf = compressImage(t, DoJPEG, JPEG, 90, "../assets/mac.jpg", "jpeg", "../assets/compress/mac-q90.jpg")
	if len(buf) != jpegSize {
		t.Fatalf("compress jpeg fail")
	}
	buf = compressImage(t, DoPNG, JPEG, 80, "../assets/mac.jpg", "png", "../assets/compress/mac.png")
	if len(buf) != pngSize {
		t.Fatalf("convert jpeg to png fail")
	}

	buf = compressImage(t, DoGUEZLI, GUETZLI, 90, "../assets/mac.jpg", "guetzli", "../assets/compress/mac-q90.guetzli.jpg")
	if len(buf) != jpegGuezliSize {
		t.Fatalf("convert jpeg(guetzli) fail")
	}
}

func resizeImage(t *testing.T, file, dst string, imageType, quality, width, height int) []byte {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatalf("get file:%s fail", file)
		return nil
	}
	newBuf, err := DoResize(buf, imageType, quality, width, height, imageType)
	if err != nil {
		t.Fatalf("resize image fail, %s %v", file, err)
		return nil
	}
	ioutil.WriteFile(dst, newBuf, os.ModePerm)
	return newBuf
}

func TestResize(t *testing.T) {
	width := 480
	buf := resizeImage(t, "../assets/compress/mac.webp", "../assets/compress/mac-resize.webp", WEBP, 75, width, 0)
	if len(buf) != webp480Size {
		t.Fatalf("resize webp fail")
	}

	buf = resizeImage(t, "../assets/compress/mac-q90.jpg", "../assets/compress/mac-resize.jpg", JPEG, 90, width, 0)
	if len(buf) != jpeg480Size {
		t.Fatalf("resize jpg fail")
	}

	buf = resizeImage(t, "../assets/compress/mac.png", "../assets/compress/mac-resize.png", PNG, 80, width, 0)
	if len(buf) != png480Size {
		t.Fatalf("resize png fail")
	}

	buf = resizeImage(t, "../assets/compress/mac.webp", "../assets/compress/mac-resize.webp", WEBP, 0, width, 0)
	if len(buf) != webpResizeSize {
		t.Fatalf("resize webp fail")
	}
}
