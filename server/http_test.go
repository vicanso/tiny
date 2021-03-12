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

package server

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
	"github.com/vicanso/go-axios"
)

var (
	jpegBase64  = "/9j/4AAQSkZJRgABAQAAAQABAAD/2wCEAAYGBgYHBgcICAcKCwoLCg8ODAwODxYQERAREBYiFRkVFRkVIh4kHhweJB42KiYmKjY+NDI0PkxERExfWl98fKcBBgYGBgcGBwgIBwoLCgsKDw4MDA4PFhAREBEQFiIVGRUVGRUiHiQeHB4kHjYqJiYqNj40MjQ+TERETF9aX3x8p//CABEIAB4APAMBIQACEQEDEQH/xAAcAAACAgIDAAAAAAAAAAAAAAAGBwQFAAECAwj/2gAIAQEAAAAA9ORbRcH0zM4wQFV3DhtrMVDC8ZELXbs2m2XrrCppx//EABQBAQAAAAAAAAAAAAAAAAAAAAD/2gAIAQIQAAAAAAD/xAAUAQEAAAAAAAAAAAAAAAAAAAAA/9oACAEDEAAAAAAA/8QAMBAAAgICAQMCAwUJAAAAAAAAAQIDBAURAAYSIRMiMUFhFDJCUYEQFRYgJGJxkbL/2gAIAQEAAT8AvZGpQRJLLOELdvcsbyefr2A6H1PI89hJVkZcnW1GCXHqAFf8g+RyN45ER0dWVhtSDsEfmDzJ9fRYrLz1rFYvFCSjtF79MdFD3bA2R4KnlC9Wv0q9yu+4pYw678HR/MfyEct0cfaT+qqwSqPJ9SNWHj5+RzMVenqDGXGZl6kkbFnrQSuYm15IYRb7OXb2QRMhPFFOlLIlHdbGmaX3bJVgBvmLy3ViVfTo5NjXiBKxoqyOIQ2iyhhshfmDzHdRhqdQ2Kd9g6gCdYhMj/37hL6B5jsxjcmHanbjl7G06qSGU/VToj9kuEzcs7Fuprawe3SRxxI/124Hz+WudXVKNKHs3cyFgFAws2XZI1kbQ7ghTZY8PRfT7wSRGvKO9CG7Z5Qu9aJ7d65c6Bx1Cg81m5Yt1aytIlZUVXbX4e8e4jjUGmyeOn6Zglb7MPWdhBpoj4Pps/sEh5ncoTWrnAyz1reTmaG3STx2TD73tPlHJP68w747+L+nYcTXkrNHBIl5HHbJsA9wk5+vG8DnVGNhiySQRT3GE1pLVjvnADefwgJ4YfI8vzZStZr16sFaQMnl5ZWUjXj4Kp3x4eqXIImxcUWtn2SyP/0vGx+dkcSHPBY/gY0qoP8ARcvzrrDy1rWPzsEoFitLFHM/gNIT4RvgRscxnTuWaazn5LsP7xm7exACYAn3fTbY2QfzHkco2ftVKrZ7O31olft3vWxz/8QAFBEBAAAAAAAAAAAAAAAAAAAAQP/aAAgBAgEBPwAH/8QAFBEBAAAAAAAAAAAAAAAAAAAAQP/aAAgBAwEBPwAH/9k="
	jpegData, _ = base64.StdEncoding.DecodeString(jpegBase64)
)

func TestOptimImageFromURL(t *testing.T) {

	t.Run("uri is nil", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/", nil)
		c := elton.NewContext(nil, req)
		err := optimImageFromURL(c)
		assert.Equal(errURLIsNil, err)
	})

	t.Run("content type is nil", func(t *testing.T) {
		assert := assert.New(t)
		done := ins.Mock(&axios.Response{
			Headers: make(http.Header),
		})
		defer done()
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/", nil)
		c := elton.NewContext(nil, req)
		err := optimImageFromURL(c)
		assert.Equal(errContentTypeIsNil, err)
	})

	t.Run("content type is invalid", func(t *testing.T) {
		assert := assert.New(t)
		headers := make(http.Header)
		headers.Set(elton.HeaderContentType, "image")
		done := ins.Mock(&axios.Response{
			Headers: headers,
		})
		defer done()
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/", nil)
		c := elton.NewContext(nil, req)
		err := optimImageFromURL(c)
		assert.Equal(errContentTypeIsInvalid, err)
	})

	t.Run("content type is not supported", func(t *testing.T) {
		assert := assert.New(t)
		headers := make(http.Header)
		headers.Set(elton.HeaderContentType, "image/bmp")
		done := ins.Mock(&axios.Response{
			Headers: headers,
		})
		defer done()
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/", nil)
		c := elton.NewContext(nil, req)
		err := optimImageFromURL(c)
		assert.Equal(errContentTypeIsNotSupported, err)
	})

	t.Run("optim jpeg", func(t *testing.T) {
		assert := assert.New(t)
		headers := make(http.Header)
		headers.Set(elton.HeaderContentType, "image/jpeg")

		done := ins.Mock(&axios.Response{
			Headers: headers,
			Data:    jpegData,
		})
		defer done()
		resp := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/&width=30", nil)
		c := elton.NewContext(resp, req)
		err := optimImageFromURL(c)
		assert.Nil(err)
		assert.NotNil(c.BodyBuffer)
	})
}

func TestOptimImageFromData(t *testing.T) {
	t.Run("image is nil", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{}`)
		err := optimImageFromData(c)
		assert.Equal(errImageIsNil, err)
	})

	t.Run("base64 encode fail", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "a"
		}`)
		err := optimImageFromData(c)
		assert.NotNil(err)
		assert.True(strings.Contains(err.Error(), "message=illegal base64 data"))
	})

	t.Run("not support source type", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "` + jpegBase64 + `"
		}`)
		err := optimImageFromData(c)
		assert.Equal(errContentTypeIsNotSupported, err)
	})

	t.Run("optim jpeg", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "` + jpegBase64 + `",
			"source": "jpeg"
		}`)
		err := optimImageFromData(c)
		assert.Nil(err)
		assert.NotNil(c.Body)
	})
}

func TestOptimTextFromURL(t *testing.T) {
	t.Run("url is nil", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/", nil)
		c := elton.NewContext(nil, req)
		err := optimTextFromURL(c)
		assert.Equal(errURLIsNil, err)
	})

	t.Run("get data fail", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/", nil)
		done := ins.Mock(&axios.Response{
			Status: 400,
		})
		defer done()

		c := elton.NewContext(nil, req)
		err := optimTextFromURL(c)
		assert.Equal("http request fail", err.Error())
	})

	t.Run("invalid output", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/", nil)
		done := ins.Mock(&axios.Response{
			Data: []byte("abcd"),
		})
		defer done()
		c := elton.NewContext(nil, req)
		err := optimTextFromURL(c)
		assert.Equal(errOutputTypeIsInvalid, err)
	})

	t.Run("gzip text", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/&output=gzip", nil)
		done := ins.Mock(&axios.Response{
			Data: []byte("abcd"),
		})
		defer done()
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		err := optimTextFromURL(c)
		assert.Nil(err)
		assert.Equal("gzip", c.GetHeader(elton.HeaderContentEncoding))
		assert.NotNil(c.BodyBuffer)
	})

	t.Run("br text", func(t *testing.T) {
		assert := assert.New(t)
		req := httptest.NewRequest("GET", "/?url=http://www.baidu.com/&output=br", nil)
		done := ins.Mock(&axios.Response{
			Data: []byte("abcd"),
		})
		defer done()
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, req)
		err := optimTextFromURL(c)
		assert.Nil(err)
		assert.Equal("br", c.GetHeader(elton.HeaderContentEncoding))
		assert.NotNil(c.BodyBuffer)
	})
}

func TestOptimTextFromData(t *testing.T) {
	t.Run("text is nil", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{}`)
		err := optimTextFromData(c)
		assert.Equal(errTextIsNil, err)
	})

	t.Run("invalid output", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "abcd",
			"output": "zz"
		}`)
		err := optimTextFromData(c)
		assert.Equal(errOutputTypeIsInvalid, err)
	})

	t.Run("gzip", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "abce",
			"output": "gzip"
		}`)
		err := optimTextFromData(c)
		assert.Nil(err)
		assert.NotNil(c.Body)
	})

	t.Run("brotli", func(t *testing.T) {
		assert := assert.New(t)
		c := elton.NewContext(nil, nil)
		c.RequestBody = []byte(`{
			"data": "abce",
			"output": "br"
		}`)
		err := optimTextFromData(c)
		assert.Nil(err)
		assert.NotNil(c.Body)
	})
}
