package tiny

import (
	"testing"

	"github.com/google/brotli/go/cbrotli"
)

func TestGzip(t *testing.T) {
	data := "daojrej123ojojaoj"
	compressBuf, err := DoGzip([]byte(data), 0)
	if err != nil {
		t.Fatalf("gzip fail, %v", err)
	}
	raw, err := DoGunzip(compressBuf)
	if err != nil {
		t.Fatalf("gzip read all fail, %v", err)
	}
	if string(raw) != data {
		t.Fatalf("gzip fail, data is not match")
	}
}

func TestBrotli(t *testing.T) {
	data := "aojdaojeanf1231ojf2o1"
	compressBuf, err := DoBrotli([]byte(data), 0)
	if err != nil {
		t.Fatalf("brotli fail, %v", err)
	}
	raw, err := cbrotli.Decode(compressBuf)
	if err != nil {
		t.Fatalf("brotli decode fail, %v", err)
	}
	if string(raw) != data {
		t.Fatalf("brotli fail, data is not match")
	}
}
