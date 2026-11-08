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
	startServer "kialkuz/service-metrics-and-alerting/pkg/start_server"
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

		if appServer.LogsPublisher != nil {
			appServer.LogsPublisher.Shutdown()
		}
	}()

	go appServer.FileService.SaveWithInterval(ctx)
	go appServer.LogsPublisher.Save(ctx)

	newServer := serverInit.NewServer(
		appServer.Handler,
		appConfigData,
		sugarlogger,
	)
	return startServer.StartWithGracefulShutdown(newServer, appConfigData.ServerPort)
}
