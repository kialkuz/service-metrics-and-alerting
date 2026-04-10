package agent

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 1
	defaultPollInterval   = 1
)

type incomingParams struct {
	address        string `env:"ADDRESS"`
	reportInterval int    `env:"REPORT_INTERVAL"`
	pollInterval   int    `env:"POLL_INTERVAL"`
}

type Config struct {
	ReportInterval int
	PollInterval   int
	URL            string
}

var config Config

func NewConfig() (*Config, error) {
	incomingParams, err := getIncomingParams()
	if err != nil {
		return nil, err
	}

	return &Config{
		ReportInterval: incomingParams.reportInterval,
		PollInterval:   incomingParams.pollInterval,
		URL:            "http://" + incomingParams.address,
	}, nil
}

func getIncomingParams() (*incomingParams, error) {
	var cfg incomingParams
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	address := flag.String("a", defaultServerAddress, "server address")
	reportInterval := flag.Int("r", defaultReportInterval, "report interval")
	pollInterval := flag.Int("p", defaultPollInterval, "poll interval")

	flag.Parse()

	if cfg.address == "" {
		cfg.address = *address
	}

	if cfg.reportInterval == 0 {
		cfg.reportInterval = *reportInterval
	}

	if cfg.pollInterval == 0 {
		cfg.pollInterval = *pollInterval
	}

	return &cfg, nil
}
