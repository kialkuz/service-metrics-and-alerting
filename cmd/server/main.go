package main

import (
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/server"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"log"
)

func main() {
	appConfigData, err := appConfig.NewConfig()
	if err != nil {
		panic(err)
	}

	storage, err := db.NewStorage(appConfigData.DBType, appConfigData.DB.DatabaseURI)
	if err != nil {
		panic(err)
	}
	defer storage.Close()

	services := service.NewMetricsService(storage)
	handler := handler.NewMetricsHandler(services)

	log.Println("Server running on port ", appConfigData.ServerPort)
	newServer := server.NewServer(handler, appConfigData.ServerPort)
	err = newServer.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
