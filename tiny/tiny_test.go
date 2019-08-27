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
// limitations under the License.

package tiny

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	pngWidth  = 80
	pngHeight = 40
	pngBase64 = "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"
)

func getTestImage() image.Image {
	buf, err := base64.StdEncoding.DecodeString(pngBase64)
	if err != nil {
		panic(err)
	}
	img, err := png.Decode(bytes.NewBuffer(buf))
	if err != nil {
		panic(err)
	}
	return img
}

func TestImageDecode(t *testing.T) {
	originalImage := getTestImage()
	t.Run("unknown image type decode", func(t *testing.T) {
		assert := assert.New(t)
		data, err := JPEGEncode(originalImage, 10)
		assert.Nil(err)
		assert.NotNil(data)
		img, err := imageDecode(data, EncodeTypeUnknown)
		assert.Nil(err)
		assert.NotNil(img)
	})

	t.Run("webp decode", func(t *testing.T) {
		assert := assert.New(t)
		data, err := WEBPEncode(originalImage, 10)
		assert.Nil(err)
		assert.NotNil(data)
		img, err := imageDecode(data, EncodeTypeWEBP)
		assert.Nil(err)
		assert.NotNil(img)
	})

	t.Run("png decode", func(t *testing.T) {
		assert := assert.New(t)
		data, err := PNGEncode(originalImage, 10)
		assert.Nil(err)
		assert.NotNil(data)
		img, err := imageDecode(data, EncodeTypePNG)
		assert.Nil(err)
		assert.NotNil(img)
	})

	t.Run("jpeg decode", func(t *testing.T) {
		assert := assert.New(t)
		data, err := JPEGEncode(originalImage, 10)
		assert.Nil(err)
		assert.NotNil(data)
		img, err := imageDecode(data, EncodeTypeJPEG)
		assert.Nil(err)
		assert.NotNil(img)
	})
}

func TestImageOptim(t *testing.T) {
	originalData, _ := base64.StdEncoding.DecodeString(pngBase64)
	t.Run("convert to webp", func(t *testing.T) {
		assert := assert.New(t)
		img, err := ImageOptim(originalData, EncodeTypePNG, EncodeTypeWEBP, 0, 0, 0)
		assert.Nil(err)
		assert.Equal(pngWidth, img.Width)
		assert.Equal(pngHeight, img.Height)
		assert.Equal(EncodeTypeWEBP, img.Type)
	})

	t.Run("convert to jpeg", func(t *testing.T) {
		assert := assert.New(t)
		width := 40
		height := 20
		img, err := ImageOptim(originalData, EncodeTypePNG, EncodeTypeJPEG, 0, width, height)
		assert.Nil(err)
		assert.Equal(width, img.Width)
		assert.Equal(height, img.Height)
		assert.Equal(EncodeTypeJPEG, img.Type)
	})

	t.Run("convert to png", func(t *testing.T) {
		assert := assert.New(t)
		width := 20
		height := 0
		img, err := ImageOptim(originalData, EncodeTypePNG, EncodeTypePNG, 0, width, height)
		assert.Nil(err)
		assert.Equal(width, img.Width)
		assert.Equal(10, img.Height)
		assert.Equal(EncodeTypePNG, img.Type)
	})
}

func TestTextOptim(t *testing.T) {
	t.Run("gzip", func(t *testing.T) {
		assert := assert.New(t)
		info, err := TextOptim([]byte("abcd"), EncodeTypeGzip, 0)
		assert.Nil(err)
		assert.Equal(EncodeTypeGzip, info.Type)
		assert.NotNil(info.Data)
	})

	t.Run("brotli", func(t *testing.T) {
		assert := assert.New(t)
		info, err := TextOptim([]byte("abcd"), EncodeTypeBr, 0)
		assert.Nil(err)
		assert.Equal(EncodeTypeBr, info.Type)
		assert.NotNil(info.Data)
	})
}

func TestEncodeType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(Gzip, EncodeTypeGzip.String())
	assert.Equal(Br, EncodeTypeBr.String())
	assert.Equal(JPEG, EncodeTypeJPEG.String())
	assert.Equal(PNG, EncodeTypePNG.String())
	assert.Equal(WEBP, EncodeTypeWEBP.String())
}

func TestConvertToEncodeType(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(EncodeTypeUnknown, ConvertToEncodeType(""))
	assert.Equal(EncodeTypeGzip, ConvertToEncodeType(Gzip))
	assert.Equal(EncodeTypeBr, ConvertToEncodeType(Br))
	assert.Equal(EncodeTypeJPEG, ConvertToEncodeType(JPEG))
	assert.Equal(EncodeTypePNG, ConvertToEncodeType(PNG))
	assert.Equal(EncodeTypeWEBP, ConvertToEncodeType(WEBP))
}
