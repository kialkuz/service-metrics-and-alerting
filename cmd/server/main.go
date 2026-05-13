package main

import (
	"context"
	"fmt"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/file"
	"kialkuz/service-metrics-and-alerting/internal/server"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	appConfigData, err := appConfig.NewConfig()
	if err != nil {
		return fmt.Errorf("error get configuration: %s", err.Error())
	}

	pool, err := pgxpool.New(context.Background(), appConfigData.DB.DatabaseURI)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %s", err.Error())
	}
	defer pool.Close()

	dbStorage := db.NewDBStorage(pool)

	fileStorage, err := file.NewFileStorage(appConfigData)
	if err != nil {
		return fmt.Errorf("error init fileStorage: %s", err.Error())
	}

	err = RestoreMetrics(appConfigData.FileStoragePath, dbStorage)
	if err != nil {
		return fmt.Errorf("error restore metrics: %s", err.Error())
	}

	fileService := service.NewFileService(fileStorage, appConfigData.StoreInterval)

	handler := handler.NewMetricsHandler(
		service.NewMetricsService(dbStorage),
		fileService,
	)

	newServer := server.NewServer(handler, appConfigData.ServerPort)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go fileService.SaveWithInterval(ctx)

	fmt.Println("Server running on port ", appConfigData.ServerPort)
	err = newServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err.Error())
	}

	return nil
}

func RestoreMetrics(fileStoragePath string, dbStorage db.MetricsDBRepository) error {
	metrics, err := file.GetFromFile(fileStoragePath)
	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		ctx := context.Background()

		dbStorage.AddList(ctx, metrics)
	}

	return nil
}
