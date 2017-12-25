package shadow

import (
	"bytes"
	"compress/gzip"

	cbrotli "github.com/google/brotli/go/cbrotli"
)

// brotli压缩
func doBrotli(buf []byte, quality int) ([]byte, error) {
	if quality == 0 {
		quality = 9
	}
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: quality,
		LGWin:   0,
	})
}

// gzip压缩
func doGzip(buf []byte, quality int) ([]byte, error) {
	if quality == 0 {
		quality = gzip.DefaultCompression
	}
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, quality)
	defer w.Close()
	if err != nil {
		return nil, err
	}
	_, err = w.Write(buf)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
