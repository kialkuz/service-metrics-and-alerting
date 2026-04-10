package server

import (
	"flag"
	"os"

	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

const (
	defaultDBType        = "postgres"
	defaultServerAddress = "localhost:8080"
)

type incomingParams struct {
	dbType  string
	address string
}

func GetIncomingParams() (*incomingParams, error) {
	var cfg incomingParams

	cfg.dbType = env.GetEnv("DB_TYPE", "")
	if cfg.dbType == "" {
		cfg.dbType = defaultDBType
	}

	cfg.address = os.Getenv("ADDRESS")
	if cfg.address == "" {
		address := flag.String("a", defaultServerAddress, "server address")
		flag.Parse()

		cfg.address = *address
	}

	return &cfg, nil
}
