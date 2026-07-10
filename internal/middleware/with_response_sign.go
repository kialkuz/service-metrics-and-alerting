package middleware

import (
	"bytes"
	"encoding/hex"
	signService "kialkuz/service-metrics-and-alerting/internal/service/sign"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signResponseWriter struct {
	gin.ResponseWriter

	body       bytes.Buffer
	statusCode int
}

func (w *signResponseWriter) WriteHeader(code int) {
	w.statusCode = code
}

func (w *signResponseWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *signResponseWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *signResponseWriter) Status() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}

	return w.statusCode
}

func WithResponseSign(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if key == "" {
			c.Next()
			return
		}

		writer := &signResponseWriter{
			ResponseWriter: c.Writer,
		}

		c.Writer = writer

		c.Next()

		body := writer.body.Bytes()

		sign := signService.Generate(body, key)

		writer.ResponseWriter.Header().Set(
			"HashSHA256",
			hex.EncodeToString(sign),
		)

		writer.ResponseWriter.WriteHeader(writer.Status())

		writer.ResponseWriter.Write(body)
	}
}
