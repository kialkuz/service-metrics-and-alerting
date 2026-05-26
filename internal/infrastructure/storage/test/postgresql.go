package postgresql

import (
	"context"
	"fmt"
	"kialkuz/service-metrics-and-alerting/pkg/pgerrors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const maxRetries = 3

var replayInterval = map[int]int{
	1: 1,
	2: 3,
	3: 5,
}

func UniversalQuery[T any](
	ctx context.Context,
	pool *pgxpool.Pool,
	operation func(conn *pgxpool.Conn) (T, error),
) (T, error) {
	var zero T

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return zero, fmt.Errorf("не удалось получить соединение: %w", err)
	}
	defer conn.Release()

	return operation(conn)
}

result, err := UniversalQuery(ctx, pool, func(conn *pgxpool.Conn) (pgx.CommandTag, error) {
	return conn.Exec(ctx, "INSERT INTO users (name) VALUES ($1)", "John")
})

row := UniversalQuery(ctx, pool, func(conn *pgxpool.Conn) (pgx.Row, error) {
	return conn.QueryRow(ctx, "SELECT name FROM users WHERE id = $1", 1), nil
})

func NewStorage(ctx context.Context, databaseURI string) (*pgxpool.Pool, error) {
	// pool, err := pgxpool.New(ctx, databaseURI)

	classifier := pgerrors.NewPostgresErrorClassifier()
	pool, err := doQuery(ctx, func() (*pgxpool.Pool, error) {
		pool, err := pgxpool.New(ctx, databaseURI)
		return pool, err
	}, classifier)

	// pool, err := newPool(ctx, databaseURI, classifier)
	// if err == nil {
	// 	return pool, nil
	// }

	// for attempt := 1; attempt <= maxRetries; attempt++ {
	// 	timer := time.NewTimer(time.Duration(replayInterval[attempt]) * time.Second)
	// 	<-timer.C

	// 	pool, err = newPool(ctx, databaseURI, classifier)
	// 	if err == nil {
	// 		return pool, nil
	// 	}
	// }

	return nil, err
}

func doQuery(
	ctx context.Context,
	operation func() (any, error),
	classifier *pgerrors.PostgresErrorClassifier,
) (any, error) {
	pool, err := operation()
	if err != nil {
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return nil, fmt.Errorf("unexpected error: %w\n", err)
		}
	}

	return pool, nil
}

func wrapPoolFactory(f func() (*pgxpool.Pool, error)) func() (any, error) {
	return func() (any, error) {
		_, err := pool.Exec(ctx, query)
		if err != nil {
			return nil, err
		}
		return pool, nil
	}
}

func ExecuteWithRetry(ctx context.Context, pool *pgxpool.Pool, query string) error {
	classifier := pgerrors.NewPostgresErrorClassifier()
	_, err := pool.Exec(ctx, query)

	for attempt := 0; attempt < maxRetries; attempt++ {
		_, lastErr = pool.Exec(ctx, "полезный SQL запрос")

		if lastErr == nil {
			return nil
		}

		// Определяем классификацию ошибки
		classification := classifier.Classify(lastErr)

		if classification == pgerrors.NonRetriable {
			return fmt.Errorf("Непредвиденная ошибка: %w\n", lastErr)
		}
	}

	return fmt.Errorf("операция прервана после %d попыток: %w", maxRetries, lastErr)
}

func newPool(
	ctx context.Context,
	query string,
	classifier *pgerrors.PostgresErrorClassifier,
) (*pgxpool.Pool, error) {
	if strings.HasPrefix(query, "SELECT") {

	} else {

	}
	// switch true {
	// case strings.HasPrefix(query, "SELECT"):
	//     ...
	// case strings.HasPrefix(query, "Прочитай"):
	// }

	pool, err := pgxpool.New(ctx, databaseURI)
	if err != nil {
		classification := classifier.Classify(err)
		if classification == pgerrors.NonRetriable {
			return nil, fmt.Errorf("unexpected error: %w\n", err)
		}
	}

	return pool, nil
}
