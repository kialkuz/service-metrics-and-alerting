package middleware

import (
	"bytes"
	"crypto/hmac"
	"encoding/hex"
	"io"
	signService "kialkuz/service-metrics-and-alerting/internal/service/sign"
	"net/http"

	"github.com/gin-gonic/gin"
)

func WithSign(key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		signForCheck := c.GetHeader("HashSHA256")

		if key != "" && signForCheck == "" {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if signForCheck != "" {
			expectedHash, err := hex.DecodeString(signForCheck)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			body, err := io.ReadAll(c.Request.Body)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			c.Request.Body = io.NopCloser(bytes.NewReader(body))

			actualHash := signService.Generate(body, key)

			if !hmac.Equal(actualHash, expectedHash) {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
		}

		c.Next()
	}
}
