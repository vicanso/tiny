package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/brotli/go/cbrotli"
)

// brotli压缩
func doBrotli(buf []byte, quality uint32) ([]byte, error) {
	if quality == 0 {
		quality = 9
	}
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: int(quality),
		LGWin:   0,
	})
}

// gzip压缩
func doGzip(buf []byte, quality uint32) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, int(quality))
	if err != nil {
		return nil, err
	}
	_, err = w.Write(buf)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

func log(category, file string, compressBuf, buf []byte, use time.Duration) {
	fmt.Printf("%v use %v to %d(%d%%) use %v \n", file, category, len(compressBuf), len(compressBuf)*100/len(buf), use)
}

func test(file string, quality uint32) {
	buf, _ := ioutil.ReadFile("./assets/" + file)
	start := time.Now()
	brBuf, _ := doBrotli(buf, quality)
	brUse := time.Since(start)
	start = time.Now()
	gzipBuf, _ := doGzip(buf, quality)
	gzipUse := time.Since(start)
	log("brotli", file, brBuf, buf, brUse)
	log("gzip", file, gzipBuf, buf, gzipUse)
}

func main() {
	test("angular.min.js", 8)
	test("lodash.min.js", 8)
	test("github.css", 8)
}
