package main

// @Title MetricsServer API
// @Description Сервер сбора метрик.
// @Version 1.0

import (
	_ "net/http/pprof"

	"context"
	"fmt"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	serverInit "kialkuz/service-metrics-and-alerting/internal/server"
	"kialkuz/service-metrics-and-alerting/internal/server/app"
	"log"

	_ "kialkuz/service-metrics-and-alerting/swagger"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	sugarlogger := logger.Sugar()
	defer logger.Sync()

	if err := run(sugarlogger); err != nil {
		sugarlogger.Fatal(err)
	}
}

func run(sugarlogger *zap.SugaredLogger) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfigData, err := appConfig.NewConfig()
	if err != nil {
		return fmt.Errorf("error get configuration: %s", err.Error())
	}

	appServer, err := app.NewApp(ctx, sugarlogger, appConfigData)
	if err != nil {
		return err
	}

	defer func() {
		if appServer.Pool != nil {
			appServer.Pool.Close()
		}
	}()

	go appServer.FileService.SaveWithInterval(ctx)

	sugarlogger.Infoln("Server running on port ", appConfigData.ServerPort)
	newServer := serverInit.NewServer(
		appServer.Handler,
		appConfigData,
		sugarlogger,
	)
	err = newServer.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %s", err.Error())
	}

	return nil
}
