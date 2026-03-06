package main

import (
	dbConfig "kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/router"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
	"net/http"
)

func main() {
	config := dbConfig.NewConfig()
	storage, err := db.NewStorage(config.DBType, config.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	services := service.NewMetricsService(storage)
	handerMetrics := handler.NewMetricsHandler(services)

	err = http.ListenAndServe(`:8080`, router.Init(handerMetrics))
	if err != nil {
		panic(err)
	}
}
