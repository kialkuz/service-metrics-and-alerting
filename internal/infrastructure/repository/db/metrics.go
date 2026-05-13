package db

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsDBRepository interface {
	Ping(ctx context.Context) error
	Add(ctx context.Context, typeValue, name string, value float64)
	AddList(ctx context.Context, metrics []model.Metrics)
	UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string)
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
}

type MemStorage struct {
	list map[string]map[string]*model.Metrics
	pool *pgxpool.Pool
}

var id = 1

func (r *MemStorage) Ping(ctx context.Context) error {
	if err := r.pool.Ping(ctx); err != nil {
		return err
	}

	return nil
}

func (r *MemStorage) Add(ctx context.Context, metricType, name string, value float64) {
	metrics := &model.Metrics{
		ID:    id,
		MType: metricType,
		Name:  name,
	}

	if metricType == model.Counter {
		intValue := int64(value)

		metrics.Delta = &intValue
	} else {
		metrics.Value = &value
	}

	if r.list[metricType] == nil {
		r.list[metricType] = make(map[string]*model.Metrics)
	}

	r.list[metricType][name] = metrics

	id++
}

func (r *MemStorage) AddList(ctx context.Context, metrics []model.Metrics) {
	for _, metric := range metrics {
		if metric.Delta != nil {
			r.Add(ctx, metric.MType, metric.Name, float64(*metric.Delta))
		} else {
			r.Add(ctx, metric.MType, metric.Name, *metric.Value)
		}
	}
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) {
	metrics := r.list[metricType][name]
	if metricType == model.Counter {
		intValue := int64(value)

		metrics.Delta = &intValue
	} else {
		metrics.Value = &value
	}

	r.list[metricType][name] = metrics
}

func (r *MemStorage) Get(ctx context.Context, metricType, name string) (*model.Metrics, error) {
	if metric, ok := r.list[metricType][name]; ok {
		return metric, nil
	}

	return nil, errors.New("metric not found")
}

func (r *MemStorage) GetList(ctx context.Context) ([]model.Metrics, error) {
	var metrics []model.Metrics

	if len(r.list) == 0 {
		return nil, errors.New("empty metrics list")
	}

	for _, metricsByName := range r.list {
		for _, metric := range metricsByName {
			metrics = append(metrics, *metric)
		}
	}

	return metrics, nil
}
