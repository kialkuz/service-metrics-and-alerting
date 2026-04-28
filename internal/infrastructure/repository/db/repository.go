package db

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"log"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"

	_ "github.com/lib/pq"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsDBRepository interface {
	Add(ctx context.Context, typeValue, name string, value float64) error
	UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
	Close()
}

func NewDBStorage(dbType string, databaseURI string) (MetricsDBRepository, error) {
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
			delta int8,
			value float8,
			CONSTRAINT users_pkey PRIMARY KEY (id)
		);
	`)
	if err != nil {
		log.Println(err)
	}

	return &MemStorage{db: db}, nil
}
