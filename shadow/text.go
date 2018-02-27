package shadow

import (
	"bytes"
	"compress/gzip"

	"github.com/google/brotli/go/cbrotli"
)

// DoGzip gzip压缩
func DoGzip(buf []byte, quality int) ([]byte, error) {
	if quality == 0 {
		quality = gzip.DefaultCompression
	}
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, quality)
	if err != nil {
		return nil, err
	}
	_, err = w.Write(buf)
	if err != nil {
		w.Close()
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}

// DoBrotli brotli压缩
func DoBrotli(buf []byte, quality int) ([]byte, error) {
	if quality == 0 {
		quality = 9
	}
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: quality,
		LGWin:   0,
	})
}
