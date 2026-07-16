package memory

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"sync"
)

type MemStorage struct {
	mrw  sync.RWMutex
	list map[string]map[string]*model.Metrics
}

func NewMemoryStorage() *MemStorage {
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
	r.updateList(ctx, metricsForUpdate)

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
	r.mrw.Lock()
	defer r.mrw.Unlock()

	metric.ID = id

	if r.list[metric.MType] == nil {
		r.list[metric.MType] = make(map[string]*model.Metrics)
	}

	r.list[metric.MType][metric.Name] = &metric

	id++

	return nil
}

func (r *MemStorage) updateList(ctx context.Context, metrics []model.Metrics) error {
	r.mrw.Lock()
	defer r.mrw.Unlock()

	for _, metric := range metrics {
		r.list[metric.MType][metric.Name] = &metric
	}

	return nil
}

func (r *MemStorage) UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error {
	r.mrw.Lock()
	defer r.mrw.Unlock()

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
	r.mrw.RLock()
	defer r.mrw.RUnlock()

	if metric, ok := r.list[metricType][name]; ok {
		return metric, nil
	}

	return nil, errors.New("metric not found")
}

func (r *MemStorage) GetList(ctx context.Context) ([]model.Metrics, error) {
	r.mrw.RLock()
	defer r.mrw.RUnlock()

	var metrics []model.Metrics

	if len(r.list) > 0 {
		for _, metricsByName := range r.list {
			for _, metric := range metricsByName {
				metrics = append(metrics, *metric)
			}
		}
	}

	return metrics, nil
}
