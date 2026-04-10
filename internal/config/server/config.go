package server

import (
	"strings"

	"kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

type Config struct {
	DBType     string
	ServerHost string
	ServerPort string
	URL        string
	DB         db.Config
}

func NewConfig() (*Config, error) {
	env.Load()
	incomingParams, err := GetIncomingParams()
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
