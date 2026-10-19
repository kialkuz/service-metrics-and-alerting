package server

import (
	"strings"

	"kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/config/db/postgresql"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

type Config struct {
	ServerHost      string
	ServerPort      string
	DB              db.Config
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	Key             string
	AuditFile       string
	AuditUrl        string
}

func NewConfig() (*Config, error) {
	err := env.Load()
	if err != nil {
		return nil, err
	}

	config, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	addressParts := strings.Split(config.Address, ":")

	return &Config{
		ServerHost:      addressParts[0],
		ServerPort:      addressParts[1],
		DB:              *postgresql.NewConfig(config.DatabaseDSN),
		StoreInterval:   config.StoreInterval,
		FileStoragePath: config.FileStoragePath,
		Restore:         config.Restore,
		Key:             config.Key,
		AuditFile:       config.AuditFile,
		AuditUrl:        config.AuditUrl,
	}, nil
}
