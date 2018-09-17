package tiny

import (
	"bytes"
	"errors"
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

// doCmdConvert call application to convert image
func doCmdConvert(app string, buf []byte, quality int) (target []byte, err error) {
	tmpfile, err := ioutil.TempFile("", app)
	if err != nil {
		return
	}
	originalFile := tmpfile.Name()
	targetFile := originalFile + "-" + app
	// remove file after done
	defer os.Remove(originalFile)
	defer os.Remove(targetFile)
	_, err = tmpfile.Write(buf)
	if err != nil {
		return
	}
	err = tmpfile.Close()
	if err != nil {
		return
	}
	var args []string
	switch app {
	case AppGuetzli:
		args = []string{
			"--quality",
			strconv.Itoa(quality),
			originalFile,
			targetFile,
		}
	case AppPngquant:
		args = []string{
			"--quality",
			strconv.Itoa(quality),
			originalFile,
			"--output",
			targetFile,
		}
	default:
		err = errors.New("not support this application")
		return
	}
	cmd := exec.Command(app, args...)

	err = cmd.Run()
	if err != nil {
		return
	}
	target, err = ioutil.ReadFile(targetFile)
	return
}

// decode decode iamge
func decode(buf []byte, imageType int) (img image.Image, err error) {
	reader := bytes.NewReader(buf)
	switch imageType {
	default:
		img, _, err = image.Decode(reader)
	case WEBP:
		img, err = webp.Decode(reader)
	case PNG:
		img, err = png.Decode(reader)
	case JPEG:
		fallthrough
	case GUETZLI:
		img, err = jpeg.Decode(reader)
	}
	return
}

// convertToWebp convert image to webp
func convertToWebp(writer *bytes.Buffer, img image.Image, quality int) error {
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

// convertToPng convert image to png
func convertToPng(writer *bytes.Buffer, img image.Image, quality int) error {

	buf := bytes.NewBuffer(nil)
	png.Encode(buf, img)
	data, err := doCmdConvert(AppPngquant, buf.Bytes(), quality)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

// convertToJpeg convert image to jpeg
func convertToJpeg(writer *bytes.Buffer, img image.Image, quality int) error {
	return jpeg.Encode(writer, img, &jpeg.Options{
		Quality: quality,
	})
}

// convertToGuetzli convert image to guetzli
func convertToGuetzli(writer *bytes.Buffer, img image.Image, quality int) error {
	buf := bytes.NewBuffer(nil)
	jpeg.Encode(buf, img, &jpeg.Options{
		Quality: quality,
	})
	data, err := doCmdConvert(AppGuetzli, buf.Bytes(), quality)
	if err != nil {
		return err
	}
	_, err = writer.Write(data)
	return err
}

func convertImage(buf []byte, imageType, quality, outputType int) (imageBuf []byte, err error) {
	img, err := decode(buf, imageType)
	if err != nil {
		return
	}
	writer := bytes.NewBuffer(nil)
	switch outputType {
	case WEBP:
		err = convertToWebp(writer, img, quality)
	case JPEG:
		err = convertToJpeg(writer, img, quality)
	case PNG:
		err = convertToPng(writer, img, quality)
	case GUETZLI:
		err = convertToGuetzli(writer, img, quality)
	default:
		err = errors.New("not support this output type")
	}
	if err != nil {
		return
	}
	imageBuf = writer.Bytes()
	return
}

// DoWebp conver image to webp
func DoWebp(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, WEBP)
}

// DoJPEG convert image to jpeg
func DoJPEG(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, JPEG)
}

// DoPNG conver image to png
func DoPNG(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, PNG)
}

// DoGUEZLI convert image to guezli
func DoGUEZLI(buf []byte, imageType, quality int) ([]byte, error) {
	return convertImage(buf, imageType, quality, GUETZLI)
}

// DoResize resize image
func DoResize(buf []byte, imageType, quality, width, height, outputType int) ([]byte, error) {
	img, err := decode(buf, imageType)
	if err != nil {
		return nil, err
	}
	if width != 0 || height != 0 {
		// 对图片做尺寸调整（原比例不变化 ）
		img = resize.Resize(uint(width), uint(height), img, resize.NearestNeighbor)
	}

	writer := bytes.NewBuffer(nil)
	switch outputType {
	default:
		err = convertToJpeg(writer, img, quality)
	case PNG:
		err = convertToPng(writer, img, quality)
	case WEBP:
		err = convertToWebp(writer, img, quality)
	case GUETZLI:
		err = convertToGuetzli(writer, img, quality)
	}
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}

// DoClip clip image
func DoClip(buf []byte, clipType, imageType, quality, width, height, outputType int) ([]byte, error) {

	if width == 0 && height == 0 {
		return nil, errors.New("width and height can not be 0 in clip mode")
	}

	img, err := decode(buf, imageType)
	if err != nil {
		return nil, err
	}
	dx := img.Bounds().Dx()
	dy := img.Bounds().Dy()
	if width == 0 {
		width = dx
	}
	if height == 0 {
		height = dy
	}
	if width > dx {
		width = dx
	}
	if height > dy {
		height = dy
	}
	var x0, y0 int
	// image.Rect
	switch clipType {
	case ClipLeftTop:
		x0 = 0
		y0 = 0
	case ClipTopCenter:
		x0 = (dx - width) / 2
		y0 = 0
	default:
		x0 = (dx - width) / 2
		y0 = (dy - height) / 2
	}
	x1 := x0 + width
	y1 := y0 + height
	r := image.Rect(x0, y0, x1, y1)
	if rgbImg, ok := img.(*image.YCbCr); ok {
		img = rgbImg.SubImage(r)
	} else if rgbImg, ok := img.(*image.RGBA); ok {
		img = rgbImg.SubImage(r)
	} else if rgbImg, ok := img.(*image.NRGBA); ok {
		img = rgbImg.SubImage(r)
	} else {
		return nil, errors.New("decode image fail")
	}

	writer := bytes.NewBuffer(nil)
	switch outputType {
	default:
		err = convertToJpeg(writer, img, quality)
	case PNG:
		err = convertToPng(writer, img, quality)
	case WEBP:
		err = convertToWebp(writer, img, quality)
	case GUETZLI:
		err = convertToGuetzli(writer, img, quality)
	}
	if err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}
