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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWebpEncode(t *testing.T) {
	assert := assert.New(t)
	img := getTestImage()

	// 无损
	data, err := WEBPEncode(img, 0)
	assert.Nil(err)
	assert.NotEqual(0, len(data))

	data, err = WEBPEncode(img, -20)
	assert.Nil(err)
	assert.NotEqual(0, len(data))

	img, err = WebpDecode(bytes.NewBuffer(data))
	assert.Nil(err)
	assert.NotNil(img)
}
