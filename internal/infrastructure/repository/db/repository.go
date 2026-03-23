package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func NewStorage(dbType string, databaseURI string) (MetricsRepository, error) {
	db, err := sql.Open(dbType, databaseURI)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE SEQUENCE IF NOT EXISTS users_id_seq;
		CREATE TABLE IF NOT EXISTS public.metrics (
			id int4 DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
			"type" varchar(50) NOT NULL,
			"name" varchar(50) NOT NULL,
			value float8 NOT NULL,
			CONSTRAINT users_pkey PRIMARY KEY (id)
		);
	`)
	if err != nil {
		log.Println(err)
	}

	return &MemStorage{db: db}, nil
}
