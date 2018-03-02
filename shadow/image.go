package shadow

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"

	"github.com/chai2010/webp"
	"github.com/nfnt/resize"
)

func decode(buf []byte, imageType int) (image.Image, error) {
	var img image.Image
	var err error
	reader := bytes.NewReader(buf)
	switch imageType {
	default:
		img, _, err = image.Decode(reader)
	case WEBP:
		img, err = webp.Decode(reader)
	case PNG:
		img, err = png.Decode(reader)
	case JPEG:
		img, err = jpeg.Decode(reader)
	case GUETZLI:
		img, err = jpeg.Decode(reader)
	}
	return img, err
}

// 图片转换函数
type convertFn func(*bytes.Buffer, image.Image, int) error

func convertImage(buf []byte, imageType, quality, outputType int) ([]byte, error) {
	img, err := decode(buf, imageType)
	if err != nil {
		return nil, err
	}
	writer := bytes.NewBuffer(nil)
	switch outputType {
	case WEBP:
		err = doWebp(writer, img, quality)
	case JPEG:
		err = doJpeg(writer, img, quality)
	case PNG:
		err = doPng(writer, img, quality)
	case GUETZLI:
		err = doGuetzli(writer, img, quality)
	}
	// err = fn(writer, img, quality)
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

func doWebp(writer *bytes.Buffer, img image.Image, quality int) error {
	if quality == 0 {
		return webp.Encode(writer, img, &webp.Options{
			Lossless: true,
		})
	}
	return webp.Encode(writer, img, &webp.Options{
		Lossless: false,
		Quality:  float32(quality),
	})
}

func doPng(writer *bytes.Buffer, img image.Image, quality int) error {
	return png.Encode(writer, img)
}

func doJpeg(writer *bytes.Buffer, img image.Image, quality int) error {
	return jpeg.Encode(writer, img, &jpeg.Options{
		Quality: quality,
	})
}

func doGuetzli(writer *bytes.Buffer, img image.Image, quality int) error {
	buf := bytes.NewBuffer(nil)
	jpeg.Encode(buf, img, &jpeg.Options{
		Quality: quality,
	})
	tmpfile, err := ioutil.TempFile("", "guetzli")
	if err != nil {
		return err
	}
	imgFile := tmpfile.Name()
	// 删除文件
	defer os.Remove(imgFile)
	if _, err := tmpfile.Write(buf.Bytes()); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	cmd := exec.Command("guetzli", "--quality", strconv.Itoa(quality), imgFile, imgFile)

	err = cmd.Run()
	if err != nil {
		return err
	}

	data, err := ioutil.ReadFile(imgFile)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err

}

// DoWebp 将图片以webp格式输出
func DoWebp(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, WEBP)
}

// DoJPEG 将图片以jpeg格式输出
func DoJPEG(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, JPEG)
}

// DoPNG 将图片以png格式输出
func DoPNG(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, PNG)
}

// DoGUEZLI 将图片以guezli处理输出
func DoGUEZLI(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, GUETZLI)
}

// DoResize 调整图片尺寸
func DoResize(buf []byte, imageType, quality, width, height, outputType int) ([]byte, error) {
	img, err := decode(buf, imageType)
	if err != nil {
		return nil, err
	}

	// 对图片做尺寸调整（原比例不变化 ）
	img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	writer := bytes.NewBuffer(nil)
	switch outputType {
	default:
		err = jpeg.Encode(writer, img, &jpeg.Options{
			Quality: quality,
		})
	case PNG:
		err = png.Encode(writer, img)
	case WEBP:
		if quality == 0 {
			err = webp.Encode(writer, img, &webp.Options{
				Lossless: true,
			})
		} else {
			err = webp.Encode(writer, img, &webp.Options{
				Lossless: false,
				Quality:  float32(quality),
			})
		}
	}
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}
