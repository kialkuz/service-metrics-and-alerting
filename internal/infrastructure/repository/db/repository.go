package db

import (
	"kialkuz/service-metrics-and-alerting/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBStorage(pool *pgxpool.Pool) MetricsDBRepository {
	list := make(map[string]map[string]*model.Metrics)

	return &MemStorage{list: list, pool: pool}
}
