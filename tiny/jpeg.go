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
	"context"
	"image"
	"image/jpeg"
	"strconv"
)

// JPEGEncode jpeg encode
func JPEGEncode(ctx context.Context, img image.Image, quality int) (data []byte, err error) {
	if quality <= minJPEGQuality || quality > maxJPEGQuality {
		quality = defaultJEPGQuality
	}
	w := new(bytes.Buffer)
	err = jpeg.Encode(w, img, &jpeg.Options{
		Quality: quality,
	})
	if err != nil {
		return
	}
	fn := func(originalFile, targetFile string) []string {
		return []string{
			"cjpeg",
			"-quality",
			strconv.Itoa(quality),
			"-outfile",
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
