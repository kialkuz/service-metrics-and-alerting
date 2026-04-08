package db

import (
<<<<<<< iter3
	"database/sql"
	"log"
=======
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
>>>>>>> v2

	_ "github.com/lib/pq"
)

func NewStorage(dbType string, databaseURI string) (MetricsRepository, error) {
	db, err := pgxpool.Connect(context.Background(), databaseURI)
	if err != nil {
		return nil, err
	}

<<<<<<< iter3
	_, err = db.Exec(`
		CREATE SEQUENCE users_id_seq;
		CREATE TABLE IF NOT EXISTS public.metrics (
			id int4 DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
=======
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err = db.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS public.metrics (
			id SERIAL NOT NULL,
>>>>>>> v2
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
