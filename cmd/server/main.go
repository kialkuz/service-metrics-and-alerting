package main

import (
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	dbConfig "kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	appConfigData := appConfig.NewConfig()
	dbConfigData := dbConfig.NewConfig()
	storage, err := db.NewStorage(appConfigData.DBType, dbConfigData.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	services := service.NewMetricsService(storage)
	handler := handler.NewMetricsHandler(services)

	router := gin.Default()
	router.POST("/update/:type/:name/:value", handler.AddHandler)

	err = http.ListenAndServe(":"+appConfigData.ServerPort, router)
	if err != nil {
		panic(err)
	}
}
