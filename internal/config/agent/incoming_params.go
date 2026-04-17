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
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PollInterval   int    `env:"POLL_INTERVAL"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		Address:        defaultServerAddress,
		ReportInterval: defaultReportInterval,
		PollInterval:   defaultPollInterval,
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "server address")
	flag.IntVar(&cfg.ReportInterval, "r", cfg.ReportInterval, "report interval")
	flag.IntVar(&cfg.PollInterval, "p", cfg.PollInterval, "poll interval")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
