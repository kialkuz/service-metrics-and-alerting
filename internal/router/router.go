package router

import (
	"kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/middleware"
	"kialkuz/service-metrics-and-alerting/internal/server/handler"

	"github.com/gin-gonic/gin"
)

func Init(handler *handler.MetricsHandler, config *server.Config) *gin.Engine {
	router := gin.New()
	router.LoadHTMLGlob("templates/*")
	router.Use(middleware.WithLogging)
	router.Use(middleware.WithComparing)
	router.Use(middleware.WithSign(config.Key))
	router.Use(middleware.WithResponseSign(config.Key))
	router.POST("/update/:type/:name/:value", handler.AddHandler)
	router.POST("/update/", handler.UpdateHandler)
	router.POST("/updates/", handler.UpdatesHandler)
	router.POST("/value/", handler.GetMetricHandler)
	router.GET("/value/:type/:name", handler.GetMetricValueHandler)
	router.GET("/", handler.GetListHandler)
	router.GET("/ping", handler.PingDB)

	return router
}
