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
	"image"
	"io"

	"github.com/chai2010/webp"
)

// WEBPEncode webp encode
func WEBPEncode(img image.Image, quality int) (data []byte, err error) {
	buffer := new(bytes.Buffer)
	if quality == 0 {
		err = webp.Encode(buffer, img, &webp.Options{
			Lossless: true,
		})
	} else {
		if quality <= minWEBPQuality || quality > maxWEBPQuality {
			quality = defaultWEBPQuality
		}
		err = webp.Encode(buffer, img, &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		})
	}

	if err != nil {
		return
	}
	data = buffer.Bytes()
	return
}

// WebpDecode webp decode
func WebpDecode(reader io.Reader) (image.Image, error) {
	return webp.Decode(reader)
}
