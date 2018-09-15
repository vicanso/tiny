package tiny

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"

	"github.com/google/brotli/go/cbrotli"
)

// DoGunzip gunzip
func DoGunzip(buf []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

// DoGzip gzip compress
func DoGzip(buf []byte, level int) ([]byte, error) {
	var b bytes.Buffer
	if level <= 0 {
		level = gzip.DefaultCompression
	}
	w, _ := gzip.NewWriterLevel(&b, level)
	_, err := w.Write(buf)
	if err != nil {
		return nil, err
	}
	// close the write to flush
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// DoBrotli brotli compress
func DoBrotli(buf []byte, quality int) ([]byte, error) {
	if quality == 0 {
		quality = 9
	}
	return cbrotli.Encode(buf, cbrotli.WriterOptions{
		Quality: quality,
		LGWin:   0,
	})
}
