package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WithLogging(c *gin.Context) {
	var sugar zap.SugaredLogger
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	sugar = *logger.Sugar()

	start := time.Now()

	c.Next()

	sugar.Infoln(
		"uri", c.Request.URL.Path,
		"method", c.Request.Method,
		"duration", time.Since(start),
		"status", c.Writer.Status(),
		"size", c.Writer.Size(),
	)
}
