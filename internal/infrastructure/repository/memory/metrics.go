package memory

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsMemoryRepository interface {
	Add(ctx context.Context, typeValue, name string, value float64) error
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
		if metric.Delta != nil {
			err = r.Add(ctx, metric.MType, metric.Name, float64(*metric.Delta))
		} else {
			err = r.Add(ctx, metric.MType, metric.Name, *metric.Value)
		}
	}

	return err
}

func (r *MemStorage) Add(ctx context.Context, metricType, name string, value float64) error {
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

	return nil
}

func (r *MemStorage) UpdateList(ctx context.Context, metrics []model.Metrics) error {
	for index, metric := range metrics {
		metric := r.list[metric.MType][metric.Name]

		switch metric.MType {
		case model.Counter:
			metrics[index].Delta = metric.Delta
		case model.Gauge:
			metrics[index].Value = metric.Value
		}
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
