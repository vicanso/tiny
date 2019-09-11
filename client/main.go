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

package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	processBar "github.com/cheggaaa/pb/v3"
	"github.com/dustin/go-humanize"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/vicanso/tiny/pb"
	"google.golang.org/grpc"
)

type (
	// Params params for optim
	Params struct {
		SourcePath  string
		TargetPath  string
		Server      string
		Filter      string
		PNGQuality  int
		JPEGQuality int
		WEBPQuality int
	}
	// OptimParams image optim params
	OptimParams struct {
		Data       []byte
		Type       string
		SourceType string
		Quality    int
	}
)

const (
	pngExt  = "png"
	webpExt = "webp"
)

// glob get match files
func glob(path, reg string) (matches []string, err error) {
	r, err := regexp.Compile(reg)
	if err != nil {
		return
	}
	matches = make([]string, 0)
	var result *multierror.Error
	filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			result = multierror.Append(result, err)
		}
		if info.IsDir() {
			return nil
		}
		if r.MatchString(file) {
			matches = append(matches, file)
		}
		return nil
	})
	err = result.ErrorOrNil()
	return
}

func getParams() (params *Params, err error) {
	source := flag.String("source", ".", "search path")
	target := flag.String("target", "", "optim target path, new image will save to this path")
	pngQuality := flag.Int("png", 90, "the quality of png, it should be >= 0 and <= 100")
	jpegQuality := flag.Int("jpeg", 80, "the quality of jpeg, it should be >= 0 and <= 100")
	webpQuality := flag.Int("webp", 0, "the quality of webp, it should be >= 0 and <= 100")
	server := flag.String("server", "tiny.aslant.site:7002", "grpc server address")
	filter := flag.String("filter", ".(png|jpg|jpeg)$", "filter regexp for image")

	flag.Parse()

	if *target == "" {
		err = errors.New("target path can not be nil")
		return
	}

	targetPath, err := filepath.Abs(*target)
	if err != nil {
		return
	}
	sourcePath, err := filepath.Abs(*source)
	if err != nil {
		return
	}

	params = &Params{
		SourcePath:  sourcePath,
		TargetPath:  targetPath,
		Server:      *server,
		Filter:      *filter,
		PNGQuality:  *pngQuality,
		JPEGQuality: *jpegQuality,
		WEBPQuality: *webpQuality,
	}
	return
}

func optim(conn *grpc.ClientConn, params *OptimParams) (data []byte, err error) {
	client := pb.NewOptimClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	in := &pb.OptimRequest{
		Data:    params.Data,
		Quality: uint32(params.Quality),
	}
	switch params.Type {
	case pngExt:
		in.Output = pb.Type_PNG
	case webpExt:
		in.Output = pb.Type_WEBP
	default:
		in.Output = pb.Type_JPEG
	}
	switch params.SourceType {
	case pngExt:
		in.Source = pb.Type_PNG
	case webpExt:
		in.Source = pb.Type_WEBP
	default:
		in.Source = pb.Type_JPEG
	}

	reply, err := client.DoOptim(ctx, in)
	if err != nil {
		return
	}
	data = reply.Data
	return
}

func getOptimParams(file string, params *Params) (optimParams *OptimParams, err error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	extType := filepath.Ext(file)[1:]
	optimParams = &OptimParams{
		Type:       extType,
		SourceType: extType,
		Data:       buf,
	}
	switch extType {
	case pngExt:
		optimParams.Quality = params.PNGQuality
	case webpExt:
		optimParams.Quality = params.WEBPQuality
	default:
		optimParams.Quality = params.JPEGQuality
	}
	return
}

func main() {
	params, err := getParams()
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
		return
	}
	result, err := glob(params.SourcePath, params.Filter)
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
		return
	}
	conn, err := grpc.Dial(params.Server, grpc.WithInsecure())
	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
		return
	}
	defer conn.Close()

	bar := processBar.StartNew(len(result))
	// 限制最多5个并发
	limiter := make(chan bool, 5)
	wg := new(sync.WaitGroup)
	var originalSizeCount uint64
	var newSizeCount uint64
	var successCount uint32
	startedAt := time.Now()
	failList := make([]string, 0)
	mutex := new(sync.Mutex)
	// 记录失败文件
	addFail := func(file string, err error) {
		mutex.Lock()
		defer mutex.Unlock()
		failList = append(failList, file+" "+err.Error())
	}
	for _, file := range result {
		limiter <- true
		go func(file string) {
			defer func() {
				bar.Increment()
				wg.Done()
				<-limiter
			}()
			wg.Add(1)
			optimParams, err := getOptimParams(file, params)
			if err != nil {
				addFail(file, err)
				return
			}
			data, err := optim(conn, optimParams)
			if err != nil {
				addFail(file, err)
				return
			}
			// 如果压缩后的数据更大，则直接使用原数据
			if len(data) >= len(optimParams.Data) {
				data = optimParams.Data
			}
			newFile := strings.Replace(file, params.SourcePath, params.TargetPath, 1)
			// 创建目录
			os.MkdirAll(filepath.Dir(newFile), os.ModePerm)
			err = ioutil.WriteFile(newFile, data, 0666)
			if err != nil {
				addFail(file, err)
				return
			}
			atomic.AddUint64(&originalSizeCount, uint64(len(optimParams.Data)))
			atomic.AddUint64(&newSizeCount, uint64(len(data)))
			atomic.AddUint32(&successCount, 1)
		}(file)
	}
	wg.Wait()
	bar.Finish()
	template := `********************************TINY********************************
Optimize images is done, use:%s
Success(%d) Fail(%d) 
Space size reduce from %s to %s
Fails: %s
********************************TINY********************************`
	errorMessage := "nil"
	if len(failList) != 0 {
		errorMessage = strings.Join(failList, "\n")
	}
	fmt.Println(fmt.Sprintf(template,
		time.Since(startedAt).String(),
		successCount,
		uint32(len(result))-successCount,
		humanize.Bytes(originalSizeCount),
		humanize.Bytes(newSizeCount),
		errorMessage,
	))
}
