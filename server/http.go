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
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/vicanso/elton"
	bodyparser "github.com/vicanso/elton-body-parser"
	recover "github.com/vicanso/elton-recover"
	responder "github.com/vicanso/elton-responder"
	"github.com/vicanso/go-axios"
	"github.com/vicanso/hes"
	"github.com/vicanso/tiny/log"
	"github.com/vicanso/tiny/tiny"
	"go.uber.org/zap"
)

type (
	optimImageParams struct {
		Data    string `json:"data,omitempty"`
		Source  string `json:"source,omitempty"`
		Output  string `json:"output,omitempty"`
		Quality int    `json:"quality,omitempty"`
		Width   int    `json:"width,omitempty"`
		Height  int    `json:"height,omitempty"`
	}
	optimTextParams struct {
		Data    string `json:"data,omitempty"`
		Output  string `json:"output,omitempty"`
		Quality int    `json:"quality,omitempty"`
	}
	// Text optim text info
	Text struct {
		Data []byte          `json:"data,omitempty"`
		Type tiny.EncodeType `json:"type,omitempty"`
	}
)

var (
	ins = axios.NewInstance(&axios.InstanceConfig{
		// 设置10秒超时
		Timeout: 10 * time.Second,
		ResponseInterceptors: []axios.ResponseInterceptor{
			newConvertResponseToError(),
		},
	})
)

// newConvertResponseToError convert http response(4xx, 5xx) to error
func newConvertResponseToError() axios.ResponseInterceptor {
	return func(resp *axios.Response) (err error) {
		if resp.Status >= 400 {
			err = errors.New("http request fail")
		}
		return
	}
}

func getIntValue(c *elton.Context, key string) int {
	v := c.QueryParam(key)
	if v == "" {
		return 0
	}
	i, _ := strconv.Atoi(v)
	return i
}

func optimImageFromURL(c *elton.Context) (err error) {
	url := c.QueryParam("url")
	if url == "" {
		err = errURLIsNil
		return
	}
	// 获取资源文件
	resp, err := ins.Get(url)
	if err != nil {
		return
	}
	// 判断content type
	contentType := resp.Headers.Get(elton.HeaderContentType)
	if contentType == "" {
		err = errContentTypeIsNil
		return
	}
	arr := strings.Split(contentType, "/")
	if len(arr) != 2 {
		err = errContentTypeIsInvalid
		return
	}
	encodeType := tiny.ConvertToEncodeType(arr[1])
	if encodeType == tiny.EncodeTypeUnknown {
		err = errContentTypeIsNotSupported
		return
	}
	outputType := tiny.ConvertToEncodeType(c.QueryParam("output"))
	// 如果未指定或不支持类型，则按保持不变
	if outputType == tiny.EncodeTypeUnknown {
		outputType = encodeType
	}
	width := getIntValue(c, "width")
	height := getIntValue(c, "height")
	quality := getIntValue(c, "quality")
	imgInfo, err := tiny.ImageOptim(resp.Data, encodeType, outputType, quality, width, height)
	if err != nil {
		return
	}
	c.SetContentTypeByExt("." + outputType.String())
	c.BodyBuffer = bytes.NewBuffer(imgInfo.Data)
	return
}

func optimImageFromData(c *elton.Context) (err error) {
	params := &optimImageParams{}
	err = json.Unmarshal(c.RequestBody, params)
	if err != nil {
		return
	}
	if params.Data == "" {
		err = errImageIsNil
		return
	}
	// 从base64转换图片数据
	data, e := base64.StdEncoding.DecodeString(params.Data)
	if e != nil {
		err = hes.Wrap(e)
		return
	}
	encodeType := tiny.ConvertToEncodeType(params.Source)
	if encodeType == tiny.EncodeTypeUnknown {
		err = errContentTypeIsNotSupported
		return
	}
	outputType := tiny.ConvertToEncodeType(params.Output)
	if outputType == tiny.EncodeTypeUnknown {
		outputType = encodeType
	}
	imgInfo, err := tiny.ImageOptim(data, encodeType, outputType, params.Quality, params.Width, params.Height)
	if err != nil {
		return
	}

	c.Body = imgInfo
	return
}

func optimTextFromURL(c *elton.Context) (err error) {
	url := c.QueryParam("url")
	if url == "" {
		err = errURLIsNil
		return
	}
	// 获取资源文件
	resp, err := ins.Get(url)
	if err != nil {
		return
	}
	outputType := tiny.ConvertToEncodeType(c.QueryParam("output"))
	if outputType == tiny.EncodeTypeUnknown {
		err = errOutputTypeIsInvalid
		return
	}
	quality := getIntValue(c, "quality")

	info, err := tiny.TextOptim(resp.Data, outputType, quality)
	if err != nil {
		return
	}

	c.SetHeader(elton.HeaderContentType, resp.Headers.Get(elton.HeaderContentType))
	c.SetHeader(elton.HeaderContentEncoding, info.Type.String())
	c.BodyBuffer = bytes.NewBuffer(info.Data)
	return
}

func optimTextFromData(c *elton.Context) (err error) {
	params := &optimTextParams{}
	err = json.Unmarshal(c.RequestBody, params)
	if err != nil {
		return
	}
	if params.Data == "" {
		err = errTextIsNil
		return
	}
	outputType := tiny.ConvertToEncodeType(params.Output)
	if outputType == tiny.EncodeTypeUnknown {
		err = errOutputTypeIsInvalid
		return
	}
	data := []byte(params.Data)
	info, err := tiny.TextOptim(data, outputType, params.Quality)
	if err != nil {
		return
	}
	c.Body = info

	return
}

// NewHTTPServer new a http server
func NewHTTPServer(address string) error {
	logger := log.Default()
	d := elton.New()
	d.OnError(func(c *elton.Context, err error) {
		// 可以针对实际场景输出更多的日志信息
		logger.DPanic("exception",
			zap.String("ip", c.RealIP()),
			zap.String("uri", c.Request.RequestURI),
			zap.Error(err),
		)
	})

	// 捕捉panic异常，避免程序崩溃
	d.Use(recover.New())
	d.Use(responder.NewDefault())

	bodyparserConf := bodyparser.Config{
		// 限制最大1MB
		Limit: 1024 * 1024,
	}
	bodyparserConf.AddDecoder(bodyparser.NewJSONDecoder())
	d.Use(bodyparser.New(bodyparserConf))

	d.GET("/ping", func(c *elton.Context) error {
		c.BodyBuffer = bytes.NewBufferString("pong")
		return nil
	})
	d.GET("/images/optim", optimImageFromURL)
	d.POST("/images/optim", optimImageFromData)

	d.GET("/texts/optim", optimTextFromURL)
	d.POST("/texts/optim", optimTextFromData)
	return d.ListenAndServe(address)
}
