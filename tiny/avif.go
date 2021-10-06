// Copyright 2021 tree xie
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
	"context"
	"image"
	"image/png"
	"strconv"
)

func AVIFEncode(ctx context.Context, img image.Image, quality int) (data []byte, err error) {
	if quality <= minAvifQuality || quality > maxAvifQuality {
		quality = defaultAvifQuality
	}
	w := new(bytes.Buffer)
	err = png.Encode(w, img)
	if err != nil {
		return
	}
	fn := func(originalFile, targetFile string) []string {
		// cavif --quality 80 --output ./test.avif -q optim
		return []string{
			"cavif",
			"--quality",
			strconv.Itoa(quality),
			"--output",
			targetFile,
			originalFile,
		}
	}
	fileBuffer := new(bytes.Buffer)
	err = doCommandConvert(ctx, w.Bytes(), fn, fileBuffer)
	if err != nil {
		return
	}
	data = fileBuffer.Bytes()
	return
}
