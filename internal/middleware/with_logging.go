package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WithLogging(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				logger.Error(err)
			}
		}

		logger.Infoln(
			"uri", c.Request.URL.Path,
			"method", c.Request.Method,
			"duration", time.Since(start),
			"status", c.Writer.Status(),
			"size", c.Writer.Size(),
		)
	}
}
