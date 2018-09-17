package tiny

import (
	"bytes"
	"io/ioutil"
	"testing"
)

const (
	pngFile  = "../assets/ai.png"
	jpgFile  = "../assets/ai.jpeg"
	webpFile = "../assets/ai.webp"
)

func TestCmdConver(t *testing.T) {
	t.Run("not support application", func(t *testing.T) {
		buf, err := ioutil.ReadFile(pngFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		_, err = doCmdConvert("xxx", buf, 90)
		if err == nil {
			t.Fatalf("not support application should return error")
		}
	})
	t.Run("png quant", func(t *testing.T) {
		buf, err := ioutil.ReadFile(pngFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		target, err := doCmdConvert(AppPngquant, buf, 90)
		if err != nil {
			t.Fatalf("png quant fail, %v", err)
		}
		if len(target) > len(buf) {
			t.Fatalf("after png quant, the file become bigger")
		}
	})

	t.Run("jpeg guetzli", func(t *testing.T) {
		buf, err := ioutil.ReadFile(jpgFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		target, err := doCmdConvert(AppGuetzli, buf, 90)
		if err != nil {
			t.Fatalf("jpeg guetzli fail, %v", err)
		}
		if len(target) > len(buf) {
			t.Fatalf("after jpeg guetzli, the file become bigger")
		}
	})
}

func TestDecode(t *testing.T) {
	t.Run("webp", func(t *testing.T) {
		buf, err := ioutil.ReadFile(webpFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		_, err = decode(buf, WEBP)
		if err != nil {
			t.Fatalf("webp decode fail, %v", err)
		}
	})

	t.Run("png", func(t *testing.T) {
		buf, err := ioutil.ReadFile(pngFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		_, err = decode(buf, PNG)
		if err != nil {
			t.Fatalf("png decode fail, %v", err)
		}
	})

	t.Run("jpeg", func(t *testing.T) {
		buf, err := ioutil.ReadFile(jpgFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		_, err = decode(buf, JPEG)
		if err != nil {
			t.Fatalf("jpeg decode fail, %v", err)
		}
	})

	t.Run("default", func(t *testing.T) {
		buf, err := ioutil.ReadFile(jpgFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}
		_, err = decode(buf, 0)
		if err != nil {
			t.Fatalf("default decode fail, %v", err)
		}
	})
}

func TestConvertToWebp(t *testing.T) {
	buf, err := ioutil.ReadFile(pngFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	img, err := decode(buf, PNG)
	if err != nil {
		t.Fatalf("png decode fail, %v", err)
	}
	t.Run("lossless", func(t *testing.T) {
		writer := bytes.NewBuffer(nil)
		err := convertToWebp(writer, img, 0)
		if err != nil {
			t.Fatalf("conver to webp lossless fail, %v", err)
		}
		if len(writer.Bytes()) == 0 {
			t.Fatalf("conver to web lossless fail, data is nil")
		}
	})

	t.Run("quality 70", func(t *testing.T) {
		writer := bytes.NewBuffer(nil)
		err := convertToWebp(writer, img, 70)
		if err != nil {
			t.Fatalf("conver to webp fail, %v", err)
		}
		if len(writer.Bytes()) == 0 {
			t.Fatalf("conver to web fail, data is nil")
		}
	})
}

func TestConvertToPng(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	img, err := decode(buf, JPEG)
	if err != nil {
		t.Fatalf("decode image fail, %v", err)
	}
	writer := bytes.NewBuffer(nil)
	err = convertToPng(writer, img, 70)
	if err != nil {
		t.Fatalf("conver to png fail, %v", err)
	}
	if len(writer.Bytes()) == 0 {
		t.Fatalf("convert to png fail, data is nil")
	}
}

func TestConvertToJpeg(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	img, err := decode(buf, JPEG)
	if err != nil {
		t.Fatalf("decode image fail, %v", err)
	}
	writer := bytes.NewBuffer(nil)
	err = convertToJpeg(writer, img, 70)
	if err != nil {
		t.Fatalf("convert to jpeg fail, %v", err)
	}
	if len(writer.Bytes()) == 0 {
		t.Fatalf("convert to jpeg fail, data is nil")
	}
}

func TestConvertToGuetzli(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	img, err := decode(buf, JPEG)
	if err != nil {
		t.Fatalf("decode image fail, %v", err)
	}
	writer := bytes.NewBuffer(nil)
	err = convertToGuetzli(writer, img, 90)
	if err != nil {
		t.Fatalf("convert to guetzli fail, %v", err)
	}
	if len(writer.Bytes()) == 0 {
		t.Fatalf("convert to guetzli fail, data is nil")
	}
}

func TestDoConvert(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	quality := 90
	imageType := JPEG
	t.Run("do return error", func(t *testing.T) {
		_, err := convertImage(buf, imageType, quality, -1)
		if err == nil {
			t.Fatalf("unknown output type should return error")
		}
	})
	t.Run("do webp", func(t *testing.T) {
		imageBuf, err := DoWebp(buf, imageType, quality)
		if err != nil {
			t.Fatalf("do webp fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("do webp fail, data is nil")
		}
	})

	t.Run("do jpeg", func(t *testing.T) {
		imageBuf, err := DoJPEG(buf, imageType, quality)
		if err != nil {
			t.Fatalf("do jpeg fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("do jpeg fail, data is nil")
		}
	})

	t.Run("do png", func(t *testing.T) {
		imageBuf, err := DoPNG(buf, imageType, quality)
		if err != nil {
			t.Fatalf("do png fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("do png fail, data is nil")
		}
	})

	t.Run("do guezli", func(t *testing.T) {
		imageBuf, err := DoGUEZLI(buf, imageType, quality)
		if err != nil {
			t.Fatalf("do guezli fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("do guezli fail, data is nil")
		}
	})
}

func TestDoResize(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	quality := 90
	imageType := JPEG
	width := 10
	height := 0
	t.Run("resize decode fail", func(t *testing.T) {
		_, err := DoResize(nil, 0, 0, 0, 0, 0)
		if err == nil {
			t.Fatalf("resize nil should return error")
		}
	})
	t.Run("resize to jpeg", func(t *testing.T) {
		imageBuf, err := DoResize(buf, imageType, quality, width, height, JPEG)
		if err != nil {
			t.Fatalf("resize to jpeg fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("resize to jpeg fail, data is nil")
		}
	})

	t.Run("resize to png", func(t *testing.T) {
		pngBuf, err := ioutil.ReadFile(pngFile)
		if err != nil {
			t.Fatalf("read file fail, %v", err)
		}

		imageBuf, err := DoResize(pngBuf, PNG, quality, width, height, PNG)
		if err != nil {
			t.Fatalf("resize to png fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("resize to png fail, data is nil")
		}
	})

	t.Run("resize to webp", func(t *testing.T) {
		imageBuf, err := DoResize(buf, imageType, quality, width, height, WEBP)
		if err != nil {
			t.Fatalf("resize to webp fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("resize to webp fail, data is nil")
		}
	})

	t.Run("resize to guetzli", func(t *testing.T) {
		imageBuf, err := DoResize(buf, imageType, quality, width, height, GUETZLI)
		if err != nil {
			t.Fatalf("resize to guetzli fail, %v", err)
		}
		if len(imageBuf) == 0 {
			t.Fatalf("resize to guetzli fail, data is nil")
		}
	})
}

func TestClip(t *testing.T) {
	buf, err := ioutil.ReadFile(jpgFile)
	if err != nil {
		t.Fatalf("read file fail, %v", err)
	}
	quality := 90
	imageType := JPEG
	width := 10
	height := 10
	t.Run("width and height both 0", func(t *testing.T) {
		_, err := DoClip(buf, ClipCenter, imageType, quality, 0, 0, WEBP)
		if err == nil {
			t.Fatalf("should return error while width and height both 0")
		}
	})

	t.Run("clip nil data", func(t *testing.T) {
		_, err := DoClip(nil, ClipCenter, imageType, quality, width, height, WEBP)
		if err == nil {
			t.Fatalf("clip nil data should return error")
		}
	})

	t.Run("clip to jpeg", func(t *testing.T) {
		_, err := DoClip(buf, ClipCenter, imageType, quality, 500, 500, JPEG)
		if err != nil {
			t.Fatalf("clip to jpeg fail, %v", err)
		}
	})

	t.Run("clip to png", func(t *testing.T) {
		_, err := DoClip(buf, ClipLeftTop, imageType, quality, width, height, PNG)
		if err != nil {
			t.Fatalf("clip to png fail, %v", err)
		}
	})

	t.Run("clip to webp", func(t *testing.T) {
		_, err := DoClip(buf, ClipCenter, imageType, quality, 0, height, WEBP)
		if err != nil {
			t.Fatalf("clip to webp fail, %v", err)
		}
	})

	t.Run("clip to guetzli", func(t *testing.T) {
		_, err := DoClip(buf, ClipCenter, imageType, quality, width, 0, GUETZLI)
		if err != nil {
			t.Fatalf("clip to guetzli fail, %v", err)
		}
	})

	t.Run("clip left top", func(t *testing.T) {
		_, err := DoClip(buf, ClipLeftTop, imageType, quality, width, height, PNG)
		if err != nil {
			t.Fatalf("clip left top to png fail, %v", err)
		}
	})

	t.Run("clip top center", func(t *testing.T) {
		_, err := DoClip(buf, ClipTopCenter, imageType, quality, width, height, PNG)
		if err != nil {
			t.Fatalf("clip top center to png fail, %v", err)
		}
	})

}
