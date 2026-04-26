package middleware

import (
	"kialkuz/service-metrics-and-alerting/internal/service/compress"
	"strings"

	"github.com/gin-gonic/gin"
)

func WithComparing(c *gin.Context) {
	acceptEncoding := c.Request.Header.Get("Content-Encoding")
	supportsGzip := strings.Contains(acceptEncoding, "gzip")

	if supportsGzip {
		unpackedBody, err := compress.ReadGzip(c.Request.Body)
		if err != nil {
			panic(err)
		}

		c.Request.Body = unpackedBody
	}

	contentEncoding := c.Request.Header.Get("Accept-Encoding")
	sendsGzip := strings.Contains(contentEncoding, "gzip")
	if sendsGzip {
		gz := compress.NewCompressWriter(c.Writer)
		if c.Writer.Status() < 300 {
			gz.Header().Add("Content-Encoding", "gzip")
		}

		defer gz.Close()
		c.Writer = gz
	}

	c.Next()
}
