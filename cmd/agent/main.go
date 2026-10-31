package main

import (
	"context"
	"errors"
	"log"

	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/agent"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	sugarlogger := logger.Sugar()
	defer logger.Sync()

	if err := run(); err != nil {
		if errors.Is(err, pkgErrors.ErrAppStopped) {
			sugarlogger.Info(pkgErrors.ErrAppStopped)
			return
		}

		sugarlogger.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config, err := appConfig.NewConfig()
	if err != nil {
		return err
	}

	services := service.NewMetricsService(config.URL, config.Key)
	agent := agent.NewMetricsAgent(services, getMetricsConfig(config))

	err = agent.CollectAndSend(ctx)
	if err != nil {
		return err
	}

	return nil
}

func getMetricsConfig(config *appConfig.Config) agent.MetricsConfig {
	return agent.MetricsConfig{
		ReportInterval: config.ReportInterval,
		PollInterval:   config.PollInterval,
		RateLimit:      config.RateLimit,
		CollectAlgo:    config.CollectAlgo,
	}
}
