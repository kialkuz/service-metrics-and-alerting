package server

import (
	"strings"

	"kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/config/db/postgresql"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

type Config struct {
	DBType          string
	ServerHost      string
	ServerPort      string
	DB              db.Config
	StoreInterval   int
	FileStoragePath string
	Restore         bool
	Key             string
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
		DBType:          config.DBType,
		ServerHost:      addressParts[0],
		ServerPort:      addressParts[1],
		DB:              *postgresql.NewConfig(config.DatabaseDSN),
		StoreInterval:   config.StoreInterval,
		FileStoragePath: config.FileStoragePath,
		Restore:         config.Restore,
		Key:             config.Key,
	}, nil
}
