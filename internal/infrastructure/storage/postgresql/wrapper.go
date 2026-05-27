package postgresql

import (
	"context"

	"kialkuz/service-metrics-and-alerting/pkg/pgerrors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const retryAttempts = 3

type RetryConfig struct {
	Attempts int
}

type DB struct {
	pool       *pgxpool.Pool
	retry      RetryConfig
	classifier *pgerrors.PostgresErrorClassifier
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool:       pool,
		retry:      RetryConfig{Attempts: retryAttempts},
		classifier: pgerrors.NewPostgresErrorClassifier(),
	}
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := WithRetry(ctx, db.retry, db.classifier, func(ctx context.Context) (struct{}, error) {
		_, err := db.pool.Exec(ctx, sql, args...)
		return struct{}{}, err
	})

	return err
}

func (db *DB) ExecTx(ctx context.Context, tx pgx.Tx, sql string, args ...any) error {
	_, err := WithRetry(ctx, db.retry, db.classifier, func(ctx context.Context) (struct{}, error) {
		_, err := tx.Exec(ctx, sql, args...)
		return struct{}{}, err
	})

	return err
}

func (db *DB) QueryRow(
	ctx context.Context,
	scan func(pgx.Row) error,
	sql string,
	args ...any,
) error {
	_, err := WithRetry(ctx, db.retry, db.classifier, func(ctx context.Context) (struct{}, error) {
		row := db.pool.QueryRow(ctx, sql, args...)
		return struct{}{}, scan(row)
	})

	return err
}
