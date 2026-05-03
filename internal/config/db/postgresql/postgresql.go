package postgresql

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/config/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/env"
)

const (
	defaultUser     = "user"
	defaultPassword = "password"
	defaultHost     = "localhost"
	defaultPort     = "5432"
	defaultDBName   = "default_name"
)

func NewConfig() *db.Config {
	return &db.Config{
		DatabaseURI: "postgres://" + getDatabaseURI(),
	}
}

func getDatabaseURI() string {
	return fmt.Sprintf("%s:%s@%s:%s/%s?sslmode=disable",
		env.GetEnv("USER", defaultUser),
		env.GetEnv("PASSWORD", defaultPassword),
		env.GetEnv("HOST", defaultHost),
		env.GetEnv("PORT", defaultPort),
		env.GetEnv("DB_NAME", defaultDBName),
	)
}
