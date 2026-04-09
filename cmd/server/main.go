package main

import (
	"flag"
	"fmt"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	dbConfig "kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/server"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
	"os"
)

const defaultServerAddress = "localhost:8080"

func main() {
	serverAddress := os.Getenv("ADDRESS")
	if serverAddress == "" {
		address := flag.String("a", defaultServerAddress, "server address")
		flag.Parse()
		serverAddress = *address
	}
	fmt.Println(serverAddress)

	appConfigData := appConfig.NewConfig(serverAddress)
	dbConfigData := dbConfig.NewConfig()
	storage, err := db.NewStorage(appConfigData.DBType, dbConfigData.DatabaseURI)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	services := service.NewMetricsService(storage)
	handler := handler.NewMetricsHandler(services)

	log.Println("Server running on port ", appConfigData.ServerPort)
	newServer := server.NewServer(handler, appConfigData.ServerPort)
	err = newServer.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
