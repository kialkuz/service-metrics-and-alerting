package db

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemStorage struct {
	pool *pgxpool.Pool
	tx   *pgx.Tx
}

func NewDBStorage(pool *pgxpool.Pool) *MemStorage {
	return &MemStorage{pool: pool}
}

func (r *MemStorage) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MemStorage) SaveList(
	ctx context.Context,
	metricsForInsert []model.Metrics,
	metricsForUpdate []model.Metrics,
) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}

	r.tx = &tx

	if len(metricsForInsert) > 0 {
		err = r.AddList(ctx, metricsForInsert)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if len(metricsForUpdate) > 0 {
		err = r.UpdateList(ctx, metricsForUpdate)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *MemStorage) AddList(ctx context.Context, metrics []model.Metrics) error {
	var (
		values []string
		args   []any
	)

	paramIndex := 1

	for _, metric := range metrics {
		values = append(values,
			fmt.Sprintf("($%d, $%d, $%d, $%d)", paramIndex, paramIndex+1, paramIndex+2, paramIndex+3),
		)

		paramIndex += 4

		if metric.Delta != nil {
			args = append(args, metric.MType, metric.Name, nil, float64(*metric.Delta))
		} else {
			args = append(args, metric.MType, metric.Name, *metric.Value, nil)
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO metrics (type, name, value, delta)
		VALUES %s
	`, strings.Join(values, ","))

	tx := *r.tx

	_, err := tx.Exec(ctx, query, args...)

	return err
}

func (r *MemStorage) UpdateList(ctx context.Context, metrics []model.Metrics) error {
	for _, metric := range metrics {
		var err error

		tx := *r.tx

		switch metric.MType {
		case model.Counter:
			_, err = tx.Exec(
				ctx,
				"UPDATE metrics SET delta = $1 WHERE id = $2",
				*metric.Delta,
				metric.ID,
			)
		case model.Gauge:
			_, err = tx.Exec(
				ctx,
				"UPDATE metrics SET value = $1 WHERE id = $2",
				*metric.Value,
				metric.ID,
			)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MemStorage) Add(ctx context.Context, metricType, name string, value float64) error {
	var err error

	if metricType == model.Counter {
		_, err = r.pool.Exec(
			ctx,
			"INSERT INTO metrics (type, name, delta) VALUES ($1, $2, $3)",
			metricType,
			name,
			int64(value),
		)
	} else {
		_, err = r.pool.Exec(
			ctx,
			"INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)",
			metricType,
			name,
			value,
		)
	}

	return err
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error {
	var err error

	if metricType == model.Counter {
		_, err = r.pool.Exec(
			ctx,
			"UPDATE metrics SET delta = $1 WHERE type = $2 AND name = $3",
			int64(value),
			metricType,
			name,
		)
	} else {
		_, err = r.pool.Exec(
			ctx,
			"UPDATE metrics SET value = $1 WHERE type = $2 AND name = $3",
			value,
			metricType,
			name,
		)
	}

	return err
}

func (r *MemStorage) Get(ctx context.Context, metricType, name string) (*model.Metrics, error) {
	metric := &model.Metrics{}

	err := r.pool.QueryRow(ctx, "SELECT * FROM metrics WHERE type = $1 AND name = $2", metricType, name).
		Scan(
			&metric.ID,
			&metric.MType,
			&metric.Name,
			&metric.Value,
			&metric.Delta,
		)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return metric, nil
}

func (r *MemStorage) GetList(ctx context.Context) ([]model.Metrics, error) {
	rows, err := r.pool.Query(ctx, "SELECT * FROM metrics")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []model.Metrics

	for rows.Next() {
		var metric model.Metrics

		if err := rows.Scan(&metric.ID, &metric.MType, &metric.Name, &metric.Value, &metric.Delta); err != nil {
			return nil, err
		}

		metrics = append(metrics, metric)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return metrics, nil
}
