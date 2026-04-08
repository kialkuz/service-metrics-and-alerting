package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/lib/pq"
)

func NewStorage(dbType string, databaseURI string) (MetricsRepository, error) {
	db, err := pgxpool.Connect(context.Background(), databaseURI)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS public.metrics (
			id SERIAL NOT NULL,
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
