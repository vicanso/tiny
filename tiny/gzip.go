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
	"compress/gzip"
)

// GzipEncode gzip compress
func GzipEncode(buf []byte, quality int) ([]byte, error) {
	if quality <= 0 || quality > gzip.BestCompression {
		quality = defaultGzipQuality
	}
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, quality)
	if err != nil {
		return nil, err
	}

	_, err = w.Write(buf)
	if err != nil {
		return nil, err
	}
	// close the write to flush
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
