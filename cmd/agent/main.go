package main

import (
	"log"

	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	"kialkuz/service-metrics-and-alerting/internal/service"
)

func main() {
	arguments, err := service.GetArguments()
	if err != nil {
		log.Fatal(err)
	}

	appConfig.NewConfig(arguments.Address)
	services := service.NewMetricsService(nil)
	agent := agent.NewMetricsAgent(services)

	agent.Collect(arguments.ReportInterval, arguments.PollInterval)
}
