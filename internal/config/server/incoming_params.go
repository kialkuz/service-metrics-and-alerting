package server

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	defaultDBType          = "postgres"
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "./storage"
	defaultRestore         = false
	defaultDatabaseDSN     = ""
)

type incomingParams struct {
	DBType          string `env:"DB_TYPE"`
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		DBType:          defaultDBType,
		Address:         defaultServerAddress,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		DatabaseDSN:     defaultDatabaseDSN,
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database dsn")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't parse config from os: %s ", err)
	}

	return cfg, nil
}
