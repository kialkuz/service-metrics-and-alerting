package middleware

import (
	"bytes"
	"encoding/hex"
	signService "kialkuz/service-metrics-and-alerting/internal/service/sign"

	"github.com/gin-gonic/gin"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w responseWriter) Write(data []byte) (int, error) {
	n, err := w.ResponseWriter.Write(data)
	if _, writeErr := w.body.Write(data); writeErr != nil && err == nil {
		err = writeErr
	}
	return n, err
}

func WithResponseSign(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.Next()
			return
		}

		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           new(bytes.Buffer),
		}

		c.Writer = writer

		c.Next()

		body := writer.body.Bytes()

		sign := signService.Generate(body, key)

		writer.ResponseWriter.Header().Set(
			"HashSHA256",
			hex.EncodeToString(sign),
		)
	}
}
