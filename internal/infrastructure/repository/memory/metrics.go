package memory

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsMemoryRepository interface {
	Add(ctx context.Context, metric model.Metrics) error
	AddList(ctx context.Context, metrics []model.Metrics) error
	SaveList(ctx context.Context, metricsForInsert []model.Metrics, metricsForUpdate []model.Metrics) error
	UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
}

type MemStorage struct {
	list map[string]map[string]*model.Metrics
}

func NewMemoryStorage() MetricsMemoryRepository {
	list := make(map[string]map[string]*model.Metrics)

	return &MemStorage{list: list}
}

var id = 1

func (r *MemStorage) SaveList(
	ctx context.Context,
	metricsForInsert []model.Metrics,
	metricsForUpdate []model.Metrics,
) error {
	r.AddList(ctx, metricsForInsert)
	r.UpdateList(ctx, metricsForUpdate)

	return nil
}

func (r *MemStorage) AddList(ctx context.Context, metrics []model.Metrics) error {
	var err error

	for _, metric := range metrics {
		err = r.Add(ctx, metric)
	}

	return err
}

func (r *MemStorage) Add(ctx context.Context, metric model.Metrics) error {
	metric.ID = id

	if r.list[metric.MType] == nil {
		r.list[metric.MType] = make(map[string]*model.Metrics)
	}

	r.list[metric.MType][metric.Name] = &metric

	id++

	return nil
}

func (r *MemStorage) UpdateList(ctx context.Context, metrics []model.Metrics) error {
	for _, metric := range metrics {
		r.list[metric.MType][metric.Name] = &metric
	}

	return nil
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error {
	metrics := r.list[metricType][name]
	if metricType == model.Counter {
		intValue := int64(value)

		metrics.Delta = &intValue
	} else {
		metrics.Value = &value
	}

	r.list[metricType][name] = metrics

	return nil
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
