package db

import (
	"database/sql"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"log"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsRepository interface {
	Add(metrics model.Metrics) error
	Update(value float64, id int) error
	Get(metricType, name string) (*model.Metrics, error)
	GetList() ([]model.Metrics, error)
	Close()
}

type MemStorage struct {
	db   *sql.DB
	list []model.Metrics
}

func (r *MemStorage) Close() {
	r.db.Close()
}

func (r *MemStorage) Add(metrics model.Metrics) error {
	_, err := r.db.Exec(`INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)`,
		metrics.MType,
		metrics.Name,
		*metrics.Value,
	)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("error add metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) Update(value float64, id int) error {
	_, err := r.db.Exec(`UPDATE metrics SET value = $1 WHERE id = $2`, value, id)
	if err != nil {
		return fmt.Errorf("error update metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) Get(metricType, name string) (*model.Metrics, error) {
	var metrics model.Metrics
	row := r.db.QueryRow(`SELECT id, type, name, value FROM metrics WHERE type=$1 AND name=$2`, metricType, name)
	err := row.Scan(&metrics.ID, &metrics.MType, &metrics.Name, &metrics.Value)
	if err != nil {
		log.Println(err)
		return nil, errors.New(pkgErrors.ErrNotFound.Error())
	}
	return &metrics, nil
}

func (r *MemStorage) GetList() ([]model.Metrics, error) {
	var metrics []model.Metrics
	rows, err := r.db.Query(`SELECT id, type, name, value FROM metrics`)
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
