package main

import (
	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	"kialkuz/service-metrics-and-alerting/internal/service"
)

func main() {
	appConfig.NewConfig()
	services := service.NewMetricsService(nil)
	agent := agent.NewMetricsAgent(services)

	agent.Collect()
}
