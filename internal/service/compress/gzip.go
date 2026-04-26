package compress

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

type gzipWriter struct {
	gin.ResponseWriter
	zw *gzip.Writer
}

func (g *gzipWriter) Write(data []byte) (int, error) {
	return g.zw.Write(data)
}

// Close закрывает gzip.Writer и досылает все данные из буфера.
func (g *gzipWriter) Close() error {
	return g.zw.Close()
}

func NewCompressWriter(w gin.ResponseWriter) *gzipWriter {
	return &gzipWriter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func MakeGzip(jsonData []byte) ([]byte, error) {
	var b bytes.Buffer

	zw := gzip.NewWriter(&b)

	if _, err := zw.Write(jsonData); err != nil {
		return nil, fmt.Errorf("error writing to gzip: %v", err)
	}

	if err := zw.Close(); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func ReadGzip(data io.Reader) (*gzip.Reader, error) {
	cr, err := gzip.NewReader(data)
	if err != nil {
		return nil, fmt.Errorf("error reading from gzip: %v", err)
	}

	defer cr.Close()

	return cr, err
}
