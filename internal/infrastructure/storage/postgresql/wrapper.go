package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RetryConfig struct {
	Attempts int
}

type DB struct {
	pool  *pgxpool.Pool
	retry RetryConfig
}

func NewDB(pool *pgxpool.Pool) *DB {
	return &DB{
		pool:  pool,
		retry: RetryConfig{Attempts: retryAttempts},
	}
}

func (db *DB) Exec(ctx context.Context, sql string, args ...any) error {
	_, err := WithRetry(ctx, db.retry, func(ctx context.Context) (struct{}, error) {
		_, err := db.pool.Exec(ctx, sql, args...)
		return struct{}{}, err
	})

	return err
}

func (db *DB) ExecTx(ctx context.Context, tx pgx.Tx, sql string, args ...any) error {
	_, err := WithRetry(ctx, db.retry, func(ctx context.Context) (struct{}, error) {
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
	_, err := WithRetry(ctx, db.retry, func(ctx context.Context) (struct{}, error) {
		row := db.pool.QueryRow(ctx, sql, args...)
		return struct{}{}, scan(row)
	})

	return err
}

// func Query[T any](
// 	ctx context.Context,
// 	db *DB,
// 	scan func(pgx.Rows) (T, error),
// 	sql string,
// 	args ...any,
// ) (T, error) {
// 	return WithRetry(ctx, db.retry, func(ctx context.Context) (T, error) {
// 		rows, err := db.pool.Query(ctx, sql, args...)
// 		if err != nil {
// 			var zero T
// 			return zero, err
// 		}
// 		defer rows.Close()

// 		return scan(rows)
// 	})
// }

// func (db *DB) WithTx(
// 	ctx context.Context,
// 	fn func(tx pgx.Tx) error,
// ) error {
// 	_, err := WithRetry(ctx, db.retry, func(ctx context.Context) (struct{}, error) {
// 		tx, err := db.pool.Begin(ctx)

// 		if err != nil {
// 			return struct{}{}, err
// 		}

// 		defer func() {
// 			if err != nil {
// 				_ = tx.Rollback(ctx)
// 			}
// 		}()

// 		err = fn(tx)
// 		if err != nil {
// 			return struct{}{}, err
// 		}

// 		return struct{}{}, tx.Commit(ctx)
// 	})

// 	return err
// }
