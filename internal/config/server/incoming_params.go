package server

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

const (
	defaultServerAddress   = "localhost:8080"
	defaultStoreInterval   = 300
	defaultFileStoragePath = "./storage"
	defaultRestore         = false
	defaultDatabaseDSN     = ""
	defaultKey             = ""
	defaultAuditFile       = ""
	defaultAuditURL        = ""
)

type incomingParams struct {
	Address         string `env:"ADDRESS"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	Restore         bool   `env:"RESTORE"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	AuditFile       string `env:"AUDIT_FILE"`
	AuditURL        string `env:"AUDIT_URL"`
}

func GetIncomingParams() (*incomingParams, error) {
	cfg := &incomingParams{
		Address:         defaultServerAddress,
		StoreInterval:   defaultStoreInterval,
		FileStoragePath: defaultFileStoragePath,
		Restore:         defaultRestore,
		DatabaseDSN:     defaultDatabaseDSN,
		Key:             defaultKey,
		AuditFile:       defaultAuditFile,
		AuditURL:        defaultAuditURL,
	}

	flag.StringVar(&cfg.Address, "a", cfg.Address, "Server address")
	flag.IntVar(&cfg.StoreInterval, "i", cfg.StoreInterval, "store interval")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	flag.BoolVar(&cfg.Restore, "r", cfg.Restore, "restore")
	flag.StringVar(&cfg.DatabaseDSN, "d", cfg.DatabaseDSN, "database dsn")
	flag.StringVar(&cfg.Key, "k", cfg.Key, "query sign")
	flag.StringVar(&cfg.AuditFile, "audit-file", cfg.AuditFile, "audit file")
	flag.StringVar(&cfg.AuditURL, "audit-url", cfg.AuditURL, "audit Url")
	flag.Parse()

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("can't parse config from os: %s ", err)
	}

	return cfg, nil
}
