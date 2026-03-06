package db

import "kialkuz/service-metrics-and-alerting/internal/infrastructure/env"

type Config struct {
	DBType      string
	Port        string
	DatabaseURI string
}

func NewConfig() *Config {
	env.Load()

	port := env.GetEnv("PORT", "8080")

	return &Config{
		DBType:      env.GetEnv("DB_TYPE", "postgres"),
		Port:        port,
		DatabaseURI: "postgres://" + getDatabaseURI(port),
	}
}

func getDatabaseURI(port string) string {
	return env.GetEnv("USER", "user") +
		":" + env.GetEnv("PASSWORD", "password") +
		"@" + env.GetEnv("HOST", "localhost") +
		":" + port +
		"/" + env.GetEnv("DB_NAME", "default_name")
}
