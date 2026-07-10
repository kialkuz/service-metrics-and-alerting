package main

import (
	"context"
	"fmt"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	serverInit "kialkuz/service-metrics-and-alerting/internal/server"
	"kialkuz/service-metrics-and-alerting/internal/server/app"
	"log"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfigData, err := appConfig.NewConfig()
	if err != nil {
		return fmt.Errorf("error get configuration: %s", err.Error())
	}

	appServer, err := app.NewApp(ctx, appConfigData)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if appServer.Pool != nil {
			appServer.Pool.Close()
		}
	}()

	go appServer.FileService.SaveWithInterval(ctx)

	fmt.Println("Server running on port ", appConfigData.ServerPort)
	newServer := serverInit.NewServer(appServer.Handler, appConfigData)
	err = newServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err.Error())
	}

	return nil
}
