package db

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"log"
<<<<<<< iter3
=======

	"github.com/jackc/pgx/v4/pgxpool"
>>>>>>> v2
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsRepository interface {
	Add(ctx context.Context, metrics model.Metrics) error
	Update(ctx context.Context, value float64, id int) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
	Close()
}

type MemStorage struct {
	db   *pgxpool.Pool
	list []model.Metrics
}

func (r *MemStorage) Close() {
	r.db.Close()
}

func (r *MemStorage) Add(ctx context.Context, metrics model.Metrics) error {
	_, err := r.db.Exec(ctx, `INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)`,
		metrics.MType,
		metrics.Name,
		*metrics.Value,
	)
	if err != nil {
<<<<<<< iter3
		log.Println("555555555555555555555")
=======
>>>>>>> v2
		log.Println(err)
		return fmt.Errorf("error add metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) Update(ctx context.Context, value float64, id int) error {
	_, err := r.db.Exec(ctx, `UPDATE metrics SET value = $1 WHERE id = $2`, value, id)
	if err != nil {
		return fmt.Errorf("error update metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) Get(ctx context.Context, metricType, name string) (*model.Metrics, error) {
	var metrics model.Metrics
<<<<<<< iter3
	row := r.db.QueryRow(`SELECT id, type, name, value FROM metrics WHERE type=$1 AND name=$2`, metricType, name)
=======
	row := r.db.QueryRow(ctx, `SELECT id, type, name, value FROM metrics WHERE type=$1 AND name=$2`, metricType, name)
>>>>>>> v2
	err := row.Scan(&metrics.ID, &metrics.MType, &metrics.Name, &metrics.Value)
	if err != nil {
		log.Println(err)
		return nil, errors.New(pkgErrors.ErrNotFound.Error())
	}
	return &metrics, nil
}

func (r *MemStorage) GetList(ctx context.Context) ([]model.Metrics, error) {
	var metrics []model.Metrics
	rows, err := r.db.Query(ctx, `SELECT id, type, name, value FROM metrics`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric model.Metrics

		err := rows.Scan(&metric.ID, &metric.MType, &metric.Name, &metric.Value)
		if err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return metrics, nil
}
