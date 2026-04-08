package main

import (
	"flag"
	"kialkuz/service-metrics-and-alerting/internal/agent"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
	"kialkuz/service-metrics-and-alerting/internal/service"
)

const (
	pollInterval   = 1
	reportInterval = 1
	host           = "localhost"
	port           = "8080"
)

func main() {
	serverAddress := flag.String(
		"a",
		env.GetEnv("SERVER_HOST", host)+":"+env.GetEnv("SERVER_PORT", port),
		"server address",
	)
	reportInterval := flag.Int("r", reportInterval, "report interval")
	pollInterval := flag.Int("p", pollInterval, "poll interval")

	flag.Parse()

	appConfig.NewConfig(*serverAddress)
	services := service.NewMetricsService(nil)
	agent := agent.NewMetricsAgent(services)

	agent.Collect(*reportInterval, *pollInterval)
}
