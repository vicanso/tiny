package shadow

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func compressText(t *testing.T, fn func([]byte, int) ([]byte, error), quality int, alg, dst string) []byte {
	buf, _ := ioutil.ReadFile("../assets/angular.min.js")
	newBuf, _ := fn(buf, quality)
	if len(newBuf) == 0 || len(newBuf) >= len(buf) {
		t.Fatalf(alg + " fail")
		return nil
	}
	log.Printf(alg+" compress, original:%d compress:%d", len(buf), len(newBuf))
	ioutil.WriteFile(dst, newBuf, os.ModePerm)
	return newBuf
}

func TestGzip(t *testing.T) {
	buf := compressText(t, DoGzip, 0, "gzip", "../assets/compress/anglar.min.js.gzip")
	if len(buf) != gzipSize {
		t.Fatalf("the gzip compress fail")
	}
}

func TestBrotli(t *testing.T) {
	buf := compressText(t, DoBrotli, 0, "brotli", "../assets/compress/anglar.min.js.br")
	if len(buf) != brotliSize {
		t.Fatalf("the br compress fail")
	}
}
