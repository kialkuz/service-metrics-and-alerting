package server

import (
	"strings"

	"kialkuz/service-metrics-and-alerting/internal/config/db"
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
}

func NewConfig() (*Config, error) {
	env.Load()

	config, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	addressParts := strings.Split(config.Address, ":")

	return &Config{
		DBType:          config.DBType,
		ServerHost:      addressParts[0],
		ServerPort:      addressParts[1],
		DB:              *db.NewConfig(),
		StoreInterval:   config.StoreInterval,
		FileStoragePath: config.FileStoragePath,
		Restore:         config.Restore,
	}, nil
}
