package server

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

type MetricsRepository interface {
	Add(ctx context.Context, typeValue, name string, value float64) error
	AddList(ctx context.Context, metrics []model.Metrics) error
	UpdateByTypeAndName(ctx context.Context, value float64, metricType, name string) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
}

//go:generate go run go.uber.org/mock/mockgen -source=metrics_service.go -destination=mocks/metrics_service_mock.go -package=mocks -typed
type MetricsService struct {
	metricsRepository MetricsRepository
}

func NewMetricsService(metricsRepository MetricsRepository) *MetricsService {
	return &MetricsService{metricsRepository: metricsRepository}
}

func (s *MetricsService) Get(ctx context.Context, metricType, name string) (*model.Metrics, error) {
	item, err := s.metricsRepository.Get(ctx, metricType, name)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *MetricsService) GetList(ctx context.Context) ([]model.Metrics, error) {
	items, err := s.metricsRepository.GetList(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *MetricsService) Save(ctx context.Context, metrics model.Metrics) error {
	existMetric, _ := s.metricsRepository.Get(ctx, metrics.MType, metrics.Name)

	switch metrics.MType {
	case model.Gauge:
		if existMetric == nil {
			return s.metricsRepository.Add(ctx, metrics.MType, metrics.Name, *metrics.Value)
		} else {
			return s.UpdateByTypeAndName(ctx, *metrics.Value, metrics.MType, metrics.Name)
		}
	case model.Counter:
		if existMetric == nil {
			return s.metricsRepository.Add(ctx, metrics.MType, metrics.Name, float64(*metrics.Delta))
		} else {
			newValue := float64(*metrics.Delta + *existMetric.Delta)

			return s.UpdateByTypeAndName(ctx, newValue, metrics.MType, metrics.Name)
		}
	}

	return errors.New("unknown metric type")
}

func (s *MetricsService) UpdateByTypeAndName(ctx context.Context, newValue float64, mType, name string) error {
	return s.metricsRepository.UpdateByTypeAndName(ctx, newValue, mType, name)
}
