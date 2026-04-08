package router

import (
	"kialkuz/service-metrics-and-alerting/internal/handler"

	"github.com/gin-gonic/gin"
)

func Init(handler *handler.MetricsHandler) *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*")
	router.POST("/update/:type/:name/:value", handler.AddHandler)
	router.GET("/value/:type/:name", handler.GetMetricHandler)
	router.GET("/", handler.GetListHandler)

	return router
}
