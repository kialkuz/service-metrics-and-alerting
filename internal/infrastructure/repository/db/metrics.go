package db

import (
	"context"
	"fmt"
	dbWrapper "kialkuz/service-metrics-and-alerting/internal/infrastructure/storage/postgresql"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MemStorage struct {
	db   *dbWrapper.DB
	pool *pgxpool.Pool
}

func NewDBStorage(db *dbWrapper.DB, pool *pgxpool.Pool) *MemStorage {
	return &MemStorage{db: db, pool: pool}
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

	if len(metricsForInsert) > 0 {
		err = r.addListWithTransaction(ctx, tx, metricsForInsert)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	if len(metricsForUpdate) > 0 {
		err = r.updateListWithTransaction(ctx, tx, metricsForUpdate)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *MemStorage) addListWithTransaction(ctx context.Context, tx pgx.Tx, metrics []model.Metrics) error {
	query, args := r.buildMultiInsertQuery(metrics)

	_, err := tx.Exec(ctx, query, args...)

	return err
}

func (r *MemStorage) AddList(ctx context.Context, metrics []model.Metrics) error {
	query, args := r.buildMultiInsertQuery(metrics)

	return r.db.Exec(ctx, query, args...)
}

func (r *MemStorage) buildMultiInsertQuery(metrics []model.Metrics) (string, []any) {
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
			args = append(args, metric.MType, metric.Name, nil, *metric.Delta)
		} else {
			args = append(args, metric.MType, metric.Name, *metric.Value, nil)
		}
	}

	return fmt.Sprintf("INSERT INTO metrics (type, name, value, delta) VALUES %s", strings.Join(values, ",")), args
}

func (r *MemStorage) updateListWithTransaction(ctx context.Context, tx pgx.Tx, metrics []model.Metrics) error {
	for _, metric := range metrics {
		var err error

		switch metric.MType {
		case model.Counter:
			query := "UPDATE metrics SET delta = $1 WHERE id = $2"
			_, err = tx.Exec(ctx, query, *metric.Delta, metric.ID)
		case model.Gauge:
			query := "UPDATE metrics SET value = $1 WHERE id = $2"
			_, err = tx.Exec(ctx, query, *metric.Value, metric.ID)
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *MemStorage) Add(ctx context.Context, metric model.Metrics) error {
	var err error

	if metric.MType == model.Counter {
		err = r.db.Exec(
			ctx,
			"INSERT INTO metrics (type, name, delta) VALUES ($1, $2, $3)",
			metric.MType,
			metric.Name,
			*metric.Delta,
		)
	} else {
		err = r.db.Exec(
			ctx,
			"INSERT INTO metrics (type, name, value) VALUES ($1, $2, $3)",
			metric.MType,
			metric.Name,
			*metric.Value,
		)
	}

	return err
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error {
	var err error

	if metricType == model.Counter {
		err = r.db.Exec(
			ctx,
			"UPDATE metrics SET delta = $1 WHERE type = $2 AND name = $3",
			int64(value),
			metricType,
			name,
		)
	} else {
		err = r.db.Exec(
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

	err := r.db.QueryRow(ctx, func(row pgx.Row) error {
		return row.Scan(
			&metric.ID,
			&metric.MType,
			&metric.Name,
			&metric.Value,
			&metric.Delta,
		)
	}, "SELECT * FROM metrics WHERE type = $1 AND name = $2", metricType, name)
	if err != nil {
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
