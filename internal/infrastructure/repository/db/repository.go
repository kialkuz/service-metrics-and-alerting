package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

type dbRepository struct {
	db *sql.DB
}

func NewStorage(dbType string, databaseURI string) (*dbRepository, error) {
	db, err := sql.Open(dbType, databaseURI)
	if err != nil {
		return nil, err
	}

	return &dbRepository{db: db}, nil
}

func (r *dbRepository) Close() {
	r.db.Close()
}
