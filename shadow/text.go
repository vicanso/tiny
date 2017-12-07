package shadow

import (
	"compress/gzip"
	"bytes"

	cbrotli "github.com/google/brotli/go/cbrotli"
)

// brotli压缩
func doBrotli(buf []byte) ([]byte, error) {
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: 9,
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