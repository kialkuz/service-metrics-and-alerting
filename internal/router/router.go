package router

import (
	"kialkuz/service-metrics-and-alerting/internal/handler"

	"github.com/gin-gonic/gin"
)

func Init(handler *handler.MetricsHandler) *gin.Engine {
	router := gin.Default()
	router.POST("/update/:type/:name/:value", handler.AddHandler)

	return router
}
