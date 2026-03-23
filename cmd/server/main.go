package main

import (
	"flag"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	dbConfig "kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/router"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
	"net/http"
)

const (
	host = "localhost"
	port = "8080"
)

func main() {
	serverAddress := flag.String(
		"a",
		env.GetEnv("SERVER_HOST", host)+":"+env.GetEnv("SERVER_PORT", port),
		"server address",
	)
	flag.Parse()

	appConfigData := appConfig.NewConfig(*serverAddress)
	dbConfigData := dbConfig.NewConfig()
	storage, err := db.NewStorage(appConfigData.DBType, dbConfigData.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	services := service.NewMetricsService(storage)
	handler := handler.NewMetricsHandler(services)

	router := router.Init(handler)

	log.Println("Server running on port ", appConfigData.ServerPort)
	if err := http.ListenAndServe(":"+appConfigData.ServerPort, router); err != nil {
		log.Fatal("Failed to start server ", err)
	}
}
