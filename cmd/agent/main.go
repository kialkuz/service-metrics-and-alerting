package main

import (
	"log"

	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/agent"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	config, err := appConfig.NewConfig()
	if err != nil {
		return err
	}
	services := service.NewMetricsService(config.URL, config.Key)
	agent := agent.NewMetricsAgent(services)

	err = agent.Collect(config.ReportInterval, config.PollInterval)
	if err != nil {
		return err
	}

	return nil
}
