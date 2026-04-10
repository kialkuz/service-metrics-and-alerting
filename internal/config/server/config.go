package config

import (
	"flag"
	"strings"

	"kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

const (
	defaultDbType        = "postgres"
	defaultServerAddress = "localhost:8080"
)

type incomingParams struct {
	dbType  string
	address string
}

type Config struct {
	DBType     string
	ServerHost string
	ServerPort string
	URL        string
	DB         db.Config
}

func NewConfig() (*Config, error) {
	env.Load()
	incomingParams, err := getIncomingParams()
	if err != nil {
		return nil, err
	}

	addressParts := strings.Split(incomingParams.address, ":")

	return &Config{
		DBType:     incomingParams.dbType,
		ServerHost: addressParts[0],
		ServerPort: addressParts[1],
		URL:        "http://" + addressParts[0] + ":" + addressParts[1],
		DB:         *db.NewConfig(),
	}, nil
}

func getIncomingParams() (*incomingParams, error) {
	var cfg incomingParams

	cfg.dbType = env.GetEnv("DB_TYPE", "")
	if cfg.dbType == "" {
		cfg.dbType = defaultDbType
	}

	cfg.address = env.GetEnv("ADDRESS", "")
	if cfg.address == "" {
		address := flag.String("a", defaultServerAddress, "server address")
		flag.Parse()

		cfg.address = *address
	}

	return &cfg, nil
}
