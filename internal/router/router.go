package router

import (
	"kialkuz/service-metrics-and-alerting/internal/middleware"
	"kialkuz/service-metrics-and-alerting/internal/server/handler"

	"github.com/gin-gonic/gin"
)

func Init(handler *handler.MetricsHandler) *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*")
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithComparing)
	router.POST("/update/:type/:name/:value", handler.AddHandler)
	router.POST("/update/", handler.UpdateHandler)
	router.POST("/updates/", handler.UpdatesHandler)
	router.POST("/value/", handler.GetMetricHandler)
	router.GET("/value/:type/:name", handler.GetMetricValueHandler)
	router.GET("/", handler.GetListHandler)
	router.GET("/ping", handler.PingDB)

	return router
}
