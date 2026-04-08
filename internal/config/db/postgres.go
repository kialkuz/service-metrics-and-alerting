package db

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

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
<<<<<<< iter3
	return env.GetEnv("USER", user) +
		":" + env.GetEnv("PASSWORD", password) +
		"@" + env.GetEnv("HOST", host) +
		":" + env.GetEnv("PORT", port) +
		"/" + env.GetEnv("DB_NAME", dbName) +
		"?sslmode=disable"
=======
	return fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		env.GetEnv("USER", user),
		env.GetEnv("PASSWORD", password),
		env.GetEnv("HOST", host),
		env.GetEnv("PORT", port),
		env.GetEnv("DB_NAME", dbName),
	)
>>>>>>> v2
}
