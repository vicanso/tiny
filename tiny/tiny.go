// Copyright 2019 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License

package tiny

import (
	"bytes"
	"errors"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

const (
	defaultGzipQuality   = 6
	defauttBrotliQuality = 6
	maxBrotliQuality     = 11

	defaultZstdQuality = 2
	minZstdQuality     = 1
	maxZstdQuality     = 2

	defaultJEPGQuality = 80
	minJPEGQuality     = 0
	maxJPEGQuality     = 100

	defaultPNGQuality = 90
	minPNGQuality     = 0
	maxPNGQuality     = 100

	defaultWEBPQuality = 80
	minWEBPQuality     = 0
	maxWEBPQuality     = 100
)

// EncodeType encode type
type EncodeType int

// CropType crop type
type CropType int

const (
	// EncodeTypeUnknown unknown
	EncodeTypeUnknown EncodeType = iota
	// EncodeTypeGzip gzip
	EncodeTypeGzip
	// EncodeTypeBr br
	EncodeTypeBr
	// EncodeTypeSnappy snappy
	EncodeTypeSnappy
	// EncodeTypeLz4 lz4
	EncodeTypeLz4
	// EncodeTypeZstd zstd
	EncodeTypeZstd
	// EncodeTypeJPEG jpeg
	EncodeTypeJPEG
	// EncodeTypePNG png
	EncodeTypePNG
	// EncodeTypeWEBP webp
	EncodeTypeWEBP
)

const (
	// CropNone none crop
	CropNone CropType = iota
	// CropLeftTop crop left top
	CropLeftTop
	// CropTopCenter crop top center
	CropTopCenter
	// CropRightTop crop right top
	CropRightTop
	// CropLeftCenter crop left center
	CropLeftCenter
	// CropCenterCenter crop center center
	CropCenterCenter
	// CropRightCenter crop right center
	CropRightCenter
	// CropLeftBottom crop left bottom
	CropLeftBottom
	// CropBottomCenter crop bottom center
	CropBottomCenter
	// CropRightBottom crop right bottom
	CropRightBottom
)

const (
	// Gzip gzip
	Gzip = "gzip"
	// Br br
	Br = "br"
	// Snappy sz
	Snappy = "sz"
	// Lz4 lz4
	Lz4 = "lz4"
	// Zstd zstd
	Zstd = "zstd"
	// JPEG jpeg
	JPEG = "jpeg"
	// PNG png
	PNG = "png"
	// WEBP webp
	WEBP = "webp"
)

type (
	// Image image information
	Image struct {
		Data   []byte     `json:"data,omitempty"`
		Type   EncodeType `json:"type,omitempty"`
		Width  int        `json:"width,omitempty"`
		Height int        `json:"height,omitempty"`
	}
	// Text text information
	Text struct {
		Data []byte     `json:"data,omitempty"`
		Type EncodeType `json:"type,omitempty"`
	}
)

func (t EncodeType) String() string {
	switch t {
	default:
		return "unknown"
	case EncodeTypeGzip:
		return Gzip
	case EncodeTypeBr:
		return Br
	case EncodeTypeSnappy:
		return Snappy
	case EncodeTypeLz4:
		return Lz4
	case EncodeTypeZstd:
		return Zstd
	case EncodeTypeJPEG:
		return JPEG
	case EncodeTypePNG:
		return PNG
	case EncodeTypeWEBP:
		return WEBP
	}
}

// ConvertToEncodeType convert to encode type
func ConvertToEncodeType(t string) EncodeType {
	switch t {
	default:
		return EncodeTypeUnknown
	case Gzip:
		return EncodeTypeGzip
	case Br:
		return EncodeTypeBr
	case Snappy:
		return EncodeTypeSnappy
	case Lz4:
		return EncodeTypeLz4
	case Zstd:
		return EncodeTypeZstd
	case JPEG:
		return EncodeTypeJPEG
	case PNG:
		return EncodeTypePNG
	case WEBP:
		return EncodeTypeWEBP
	}
}

func imageDecode(buf []byte, sourceType EncodeType) (img image.Image, err error) {
	reader := bytes.NewReader(buf)
	switch sourceType {
	default:
		img, _, err = image.Decode(reader)
	case EncodeTypeWEBP:
		img, err = WebpDecode(reader)
	case EncodeTypePNG:
		img, err = png.Decode(reader)
	case EncodeTypeJPEG:
		img, err = jpeg.Decode(reader)
	}
	return
}

// ImageResize resize image
func ImageResize(img image.Image, width, height int) image.Image {
	return imaging.Resize(img, width, height, imaging.Lanczos)
}

// ImageCrop crop image
func ImageCrop(img image.Image, cropType CropType, width, height int) image.Image {
	currentWidth := img.Bounds().Dx()
	currentHeight := img.Bounds().Dy()
	if width == 0 || width > currentWidth {
		width = currentWidth
	}
	if height == 0 || height > currentHeight {
		height = currentHeight
	}
	var x0, y0, x1, y1 int
	switch cropType {
	case CropLeftTop:
		x0 = 0
		y0 = 0
	case CropTopCenter:
		x0 = (currentWidth - width) / 2
		y0 = 0
	case CropRightTop:
		x0 = (currentWidth - width)
		y0 = 0
	case CropLeftCenter:
		x0 = 0
		y0 = (currentHeight - height) / 2
	case CropCenterCenter:
		x0 = (currentWidth - width) / 2
		y0 = (currentHeight - height) / 2
	case CropRightCenter:
		x0 = currentWidth - width
		y0 = (currentHeight - height) / 2
	case CropLeftBottom:
		x0 = 0
		y0 = currentHeight - height
	case CropBottomCenter:
		x0 = (currentWidth - width) / 2
		y0 = currentHeight - height
	case CropRightBottom:
		x0 = currentWidth - width
		y0 = currentHeight - height
	default:
		// 其它的裁切类型不处理
	}
	x1 = x0 + width
	y1 = y0 + height

	rect := image.Rect(x0, y0, x1, y1)
	return imaging.Crop(img, rect)
}

// ImageOptim image optim
func ImageOptim(buf []byte, sourceType, outputType EncodeType, cropType CropType, quality, width, height int) (imgInfo *Image, err error) {
	img, err := imageDecode(buf, sourceType)
	if err != nil {
		return
	}
	if width != 0 || height != 0 {
		// 如果不需要裁剪
		if cropType == CropNone {
			img = ImageResize(img, width, height)
		} else {
			img = ImageCrop(img, cropType, width, height)
		}
	}
	var data []byte
	switch outputType {
	case EncodeTypeJPEG:
		data, err = JPEGEncode(img, quality)
	case EncodeTypePNG:
		data, err = PNGEncode(img, quality)
	case EncodeTypeWEBP:
		data, err = WEBPEncode(img, quality)
	default:
		err = errors.New("not support the output type")
	}
	if err != nil {
		return
	}
	w := img.Bounds().Dx()
	h := img.Bounds().Dy()
	imgInfo = &Image{
		Data:   data,
		Width:  w,
		Height: h,
		Type:   outputType,
	}
	return
}

// TextOptim text optim
func TextOptim(data []byte, outputType EncodeType, quality int) (info *Text, err error) {
	var buf []byte
	t := EncodeTypeUnknown
	switch outputType {
	case EncodeTypeBr:
		t = EncodeTypeBr
		buf, err = BrotliEncode(data, quality)
	case EncodeTypeSnappy:
		t = EncodeTypeSnappy
		buf, err = SnappyEncode(data)
	case EncodeTypeLz4:
		t = EncodeTypeLz4
		buf, err = Lz4Encode(data, quality)
	case EncodeTypeZstd:
		t = EncodeTypeZstd
		buf, err = ZstdEncode(data, quality)
	default:
		t = EncodeTypeGzip
		buf, err = GzipEncode(data, quality)
	}
	if err != nil {
		return
	}
	info = &Text{
		Data: buf,
		Type: t,
	}
	return
}
