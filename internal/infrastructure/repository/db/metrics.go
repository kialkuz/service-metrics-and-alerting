package db

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"

	"github.com/jackc/pgx/v4/pgxpool"
)

type MemStorage struct {
	db   *pgxpool.Pool
	list []model.Metrics
}

func (r *MemStorage) Close() {
	r.db.Close()
}

func (r *MemStorage) Add(ctx context.Context, metricType, name string, value float64) error {
	query := fmt.Sprintf(`INSERT INTO metrics (type, name, %s) VALUES ($1, $2, $3)`, r.getValueFieldName(metricType))

	_, err := r.db.Exec(ctx, query, metricType, name, value)
	if err != nil {
		return fmt.Errorf("error add metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error {
	query := fmt.Sprintf(`UPDATE metrics SET %s = $1 WHERE "type" = $2 AND name = $3`, r.getValueFieldName(metricType))
	_, err := r.db.Exec(ctx, query, value, metricType, name)
	if err != nil {
		return fmt.Errorf("error update metrics: %w", err)
	}
	return nil
}

func (r *MemStorage) getValueFieldName(metricType string) string {
	var valueField string
	switch metricType {
	case model.Counter:
		valueField = "delta"
	case model.Gauge:
		valueField = "value"
	}

	return valueField
}

func (r *MemStorage) Get(ctx context.Context, metricType, name string) (*model.Metrics, error) {
	var metrics model.Metrics
	row := r.db.QueryRow(
		ctx,
		`SELECT id, type, name, value, delta FROM metrics WHERE "type"=$1 AND name=$2`,
		metricType,
		name,
	)
	err := row.Scan(&metrics.ID, &metrics.MType, &metrics.Name, &metrics.Value, &metrics.Delta)
	if err != nil {
		return nil, errors.New(pkgErrors.ErrNotFound.Error())
	}
	return &metrics, nil
}

func (r *MemStorage) GetList(ctx context.Context) ([]model.Metrics, error) {
	var metrics []model.Metrics
	rows, err := r.db.Query(
		ctx,
		`SELECT id, type, name, value, delta FROM metrics`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var metric model.Metrics

		err := rows.Scan(&metric.ID, &metric.MType, &metric.Name, &metric.Value, &metric.Delta)
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
