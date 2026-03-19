package db

import (
	"database/sql"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/pkg/errors"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsRepository interface {
	AddMetric(metrics model.Metrics) error
	UpdateMetric(value float64, id int) error
	GetMetric(name string) (*model.Metrics, error)
	Close()
}

type MemStorage struct {
	db   *sql.DB
	list []model.Metrics
}

func (r *MemStorage) Close() {
	r.db.Close()
}

func (r *MemStorage) AddMetric(metrics model.Metrics) error {
	_, err := r.db.Exec(`INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)`,
		metrics.MType,
		metrics.Name,
		*metrics.Value,
	)
	if err != nil {
		return fmt.Errorf("error add metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) UpdateMetric(value float64, id int) error {
	_, err := r.db.Exec(`UPDATE metrics SET value = $1 WHERE id = $2`, value, id)
	if err != nil {
		return fmt.Errorf("error update metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) GetMetric(name string) (*model.Metrics, error) {
	var metrics model.Metrics
	err := r.db.QueryRow(`SELECT id, type, name, value FROM metrics WHERE name=$1`, name).
		Scan(&metrics.ID, &metrics.MType, &metrics.Name, &metrics.Value)
	if err != nil {
		return nil, errors.ErrNotFound
	}
	return &metrics, nil
}
