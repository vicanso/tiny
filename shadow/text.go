package shadow

import (
	"bytes"
	"compress/gzip"

	cbrotli "github.com/google/brotli/go/cbrotli"
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
func doGzip(buf []byte) ([]byte, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	w.Close()
	return b.Bytes(), nil
}
