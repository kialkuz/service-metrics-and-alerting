package config

import "kialkuz/service-metrics-and-alerting/internal/infrastructure/env"

const (
	dbType = "postgres"
	host   = "localhost"
	port   = "8080"
)

type Config struct {
	DBType     string
	ServerHost string
	ServerPort string
	URL        string
}

func NewConfig() *Config {
	env.Load()

	serverHost := env.GetEnv("SERVER_HOST", host)
	serverPort := env.GetEnv("SERVER_PORT", port)

	return &Config{
		DBType:     env.GetEnv("DB_TYPE", dbType),
		ServerHost: serverHost,
		ServerPort: serverPort,
		URL:        "http://" + serverHost + ":" + serverPort,
	}
}
