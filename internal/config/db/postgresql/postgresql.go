package postgresql

import (
	"kialkuz/service-metrics-and-alerting/internal/config/db"
)

func NewConfig(databaseDSN string) *db.Config {
	return &db.Config{
		DatabaseURI: databaseDSN,
	}
}
