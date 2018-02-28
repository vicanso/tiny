package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/vicanso/tiny/shadow"
)

func log(category, file string, compressBuf, buf []byte, use time.Duration) {
	fmt.Printf("%v use %v to %d(%d%%) use %v \n", file, category, len(compressBuf), len(compressBuf)*100/len(buf), use)
}

func test(file string, quality int) {
	buf, _ := ioutil.ReadFile("./assets/" + file)
	start := time.Now()
	brBuf, _ := shadow.DoBrotli(buf, quality)
	brUse := time.Since(start)
	start = time.Now()
	gzipBuf, _ := shadow.DoGzip(buf, quality)
	gzipUse := time.Since(start)
	log("brotli", file, brBuf, buf, brUse)
	log("gzip", file, gzipBuf, buf, gzipUse)
}

func main() {
	quality := flag.Int("q", 9, "quality")
	flag.Parse()
	test("angular.min.js", *quality)
	test("lodash.min.js", *quality)
	test("github.css", *quality)
}
