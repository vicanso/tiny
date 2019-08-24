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
	defauttBrotliQuality = 6
	maxBrotliQuality     = 9

	defaultJEPGQuality = 70
	minJPEGQuality     = 0
	maxJPEGQuality     = 100

	defaultPNGQuality = 80
	minPNGQuality     = 0
	maxPNGQuality     = 100

	defaultWEBPQuality = 80
	minWEBPQuality     = 0
	maxWEBPQuality     = 100
)

// EncodeType encode type
type EncodeType int

const (
	// EncodeTypeUnknown unknown
	EncodeTypeUnknown EncodeType = iota
	// EncodeTypeGzip gzip
	EncodeTypeGzip
	// EncodeTypeBr br
	EncodeTypeBr
	// EncodeTypeJPEG jpeg
	EncodeTypeJPEG
	// EncodeTypePNG png
	EncodeTypePNG
	// EncodeTypeWEBP webp
	EncodeTypeWEBP
)

type (
	// Image image information
	Image struct {
		Data  []byte     `json:"data,omitempty"`
		Type  EncodeType `json:"type,omitempty"`
		Width int        `json:"width,omitempty"`
		Heiht int        `json:"heiht,omitempty"`
	}
)

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

// ImageOptiom image optim
func ImageOptiom(buf []byte, sourceType, outputType EncodeType, quality, width, height int) (imgInfo *Image, err error) {
	img, err := imageDecode(buf, sourceType)
	if err != nil {
		return
	}
	if width != 0 || height != 0 {
		img = ImageResize(img, width, height)
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
		Data:  data,
		Width: w,
		Heiht: h,
		Type:  outputType,
	}
	return
}
