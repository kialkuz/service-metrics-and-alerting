package postgresql

import (
	"context"

	"kialkuz/service-metrics-and-alerting/internal/service/retrier"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const retryAttempts = 3

type DB struct {
	pool    *pgxpool.Pool
	retrier retrier.Retrier
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool: pool,
		retrier: retrier.NewRetrierService(
			retrier.RetryConfig{Attempts: retryAttempts},
			errorclassifier.NewPostgresErrorClassifier(),
		),
	}
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) error {
	return db.retrier.WithRetry(ctx, func(ctx context.Context) error {
		_, err := db.pool.Exec(ctx, sql, args...)
		return err
	})
}

func (db *DB) ExecTx(ctx context.Context, tx pgx.Tx, sql string, args ...any) error {
	return db.retrier.WithRetry(ctx, func(ctx context.Context) error {
		_, err := tx.Exec(ctx, sql, args...)
		return err
	})
}

func (db *DB) QueryRow(
	ctx context.Context,
	scan func(pgx.Row) error,
	sql string,
	args ...any,
) error {
	return db.retrier.WithRetry(ctx, func(ctx context.Context) error {
		row := db.pool.QueryRow(ctx, sql, args...)
		return scan(row)
	})
}
