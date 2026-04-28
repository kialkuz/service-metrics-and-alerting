package main

import (
	"context"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/file"
	"kialkuz/service-metrics-and-alerting/internal/server"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"log"
)

func main() {
	appConfigData, err := appConfig.NewConfig()
	if err != nil {
		panic(err)
	}

	dbStorage, err := db.NewDBStorage(appConfigData.DBType, appConfigData.DB.DatabaseURI)
	if err != nil {
		panic(err)
	}
	defer dbStorage.Close()

	err = RestoreMetrics(appConfigData.FileStoragePath, dbStorage)
	if err != nil {
		panic(err)
	}

	fileStorage, err := file.NewFileStorage(appConfigData.FileStoragePath, appConfigData.StoreInterval)
	if err != nil {
		panic(err)
	}
	defer fileStorage.Close()

	fileService := service.NewFileService(fileStorage, appConfigData.StoreInterval)

	handler := handler.NewMetricsHandler(
		service.NewMetricsService(dbStorage),
		fileService,
	)

	log.Println("Server running on port ", appConfigData.ServerPort)
	newServer := server.NewServer(handler, appConfigData.ServerPort)
	err = newServer.ListenAndServe()
	if err != nil {
		panic(err)
	}

	fileService.SaveWithInterval()
}

func RestoreMetrics(fileStoragePath string, dbStorage db.MetricsDBRepository) error {
	metrics, err := file.GetFromFile(fileStoragePath)
	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		ctx := context.Background()

		for _, metric := range metrics {
			if metric.Delta != nil {
				dbStorage.Add(ctx, metric.MType, metric.Name, float64(*metric.Delta))
			} else {
				dbStorage.Add(ctx, metric.MType, metric.Name, *metric.Value)
			}
		}
	}

	return nil
}
