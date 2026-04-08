package config

import (
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
	"strings"
)

const (
	dbType = "postgres"
)

type Config struct {
	DBType     string
	ServerHost string
	ServerPort string
	URL        string
}

var config Config

func NewConfig(serverAddress string) Config {
	env.Load()

	addressParts := strings.Split(serverAddress, ":")

	config = Config{
		DBType:     env.GetEnv("DB_TYPE", dbType),
		ServerHost: addressParts[0],
		ServerPort: addressParts[1],
		URL:        "http://" + addressParts[0] + ":" + addressParts[1],
	}

	return config
}

func GetURL() string {
	return config.URL
}
