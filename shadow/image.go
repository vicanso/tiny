package shadow

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

func getImage(buf []byte, width, height uint32, imageType string) (image.Image, error) {
	var img image.Image
	var err error
	reader := bytes.NewReader(buf)
	switch imageType {
	default:
		img, _, err = image.Decode(reader)
	case PNG:
		img, err = png.Decode(reader)
	case JPEG:
		img, err = jpeg.Decode(reader)
	}

	if err != nil {
		return nil, err
	}

	if width == 0 && height == 0 {
		origBounds := img.Bounds()
		width = uint32(origBounds.Dx())
	}

	// 对图片做压缩处理（原尺寸不变化 ）
	img = resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	return img, nil
}

// 对图片做webp转换压缩
func doWebp(buf []byte, width, height, quality uint32, imageType string) ([]byte, error) {
	img, err := getImage(buf, width, height, imageType)
	if err != nil {
		return nil, err
	}
	newBuf := bytes.NewBuffer(nil) //开辟一个新的空buff
	// 默认转换质量
	if quality == 0 {
		err = webp.Encode(newBuf, img, &webp.Options{
			Lossless: true,
		})
	} else {
		err = webp.Encode(newBuf, img, &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		})
	}
	if err != nil {
		return nil, err
	}
	return newBuf.Bytes(), nil
}

// 对图片做jegp转换压缩
func doJPEG(buf []byte, width, height, quality uint32, imageType string) ([]byte, error) {
	img, err := getImage(buf, width, height, imageType)
	if err != nil {
		return nil, err
	}
	newBuf := bytes.NewBuffer(nil) //开辟一个新的空buff
	err = jpeg.Encode(newBuf, img, &jpeg.Options{
		Quality: int(quality),
	})
	if err != nil {
		return nil, err
	}
	return newBuf.Bytes(), nil
}

// 对图片做png转换压缩
func doPNG(buf []byte, width, height uint32, imageType string) ([]byte, error) {
	img, err := getImage(buf, width, height, imageType)
	if err != nil {
		return nil, err
	}
	newBuf := bytes.NewBuffer(nil) //开辟一个新的空buff
	err = png.Encode(newBuf, img)
	if err != nil {
		return nil, err
	}
	return newBuf.Bytes(), nil
}
