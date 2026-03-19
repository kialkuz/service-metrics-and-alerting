package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewStorage(dbType string, databaseURI string) (MetricsRepository, error) {
	db, err := sql.Open(dbType, databaseURI)
	if err != nil {
		return nil, err
	}

	return &MemStorage{db: db}, nil
}
