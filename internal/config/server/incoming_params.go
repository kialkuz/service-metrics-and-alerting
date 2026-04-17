package server

import (
	"flag"
	"fmt"
	"log"

	"github.com/caarlos0/env/v6"
)

const (
	defaultDBType        = "postgres"
	defaultServerAddress = "localhost:8080"
)

type incomingParams struct {
	DBType  string `env:"DB_TYPE"`
	Address string `env:"ADDRESS"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		DBType:  defaultDBType,
		Address: defaultServerAddress,
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		log.Printf("Can't parse config from os: %s ", err)
		return nil, fmt.Errorf("can't parse config from os: %s ", err)
	}

	return cfg, nil
}
