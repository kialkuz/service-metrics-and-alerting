package main

import (
	"log"

	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/agent"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
)

func main() {
	config, err := appConfig.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	services := service.NewMetricsService(config.URL)
	agent := agent.NewMetricsAgent(services)

	agent.Collect(config.ReportInterval, config.PollInterval)
}
