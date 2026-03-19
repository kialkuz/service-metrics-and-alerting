package db

import "kialkuz/service-metrics-and-alerting/internal/infrastructure/env"

const (
	port     = "5432"
	user     = "user"
	password = "password"
	host     = "localhost"
	dbName   = "default_name"
)

type Config struct {
	DatabaseURI string
}

func NewConfig() *Config {
	env.Load()

	return &Config{
		DatabaseURI: "postgres://" + getDatabaseURI(),
	}
}

func getDatabaseURI() string {
	return env.GetEnv("USER", user) +
		":" + env.GetEnv("PASSWORD", password) +
		"@" + env.GetEnv("HOST", host) +
		":" + env.GetEnv("PORT", port) +
		"/" + env.GetEnv("DB_NAME", dbName)
}
