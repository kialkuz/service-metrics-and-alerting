package db

import (
	"kialkuz/service-metrics-and-alerting/internal/model"
)

func NewDBStorage(fileStoragePath string) MetricsDBRepository {
	list := make(map[string]map[string]*model.Metrics)

	return &MemStorage{list: list}
}
