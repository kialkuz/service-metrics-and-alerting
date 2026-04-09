package service

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress  = "localhost:8080"
	defaultReportInterval = 1
	defaultPollInterval   = 1
)

type Config struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func GetArguments() (*Config, error) {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	address := flag.String("a", defaultServerAddress, "server address")
	reportInterval := flag.Int("r", defaultReportInterval, "report interval")
	pollInterval := flag.Int("p", defaultPollInterval, "poll interval")

	flag.Parse()

	if cfg.Address == "" {
		cfg.Address = *address
	}

	if cfg.ReportInterval == 0 {
		cfg.ReportInterval = *reportInterval
	}

	if cfg.PollInterval == 0 {
		cfg.PollInterval = *pollInterval
	}

	return &cfg, nil
}
