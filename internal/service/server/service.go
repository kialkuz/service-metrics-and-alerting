package server

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

type MetricsServerService interface {
	Save(ctx context.Context, metrics model.Metrics) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
	UpdateByTypeAndName(ctx context.Context, newValue float64, mType, name string)
}

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks -typed
type MetricsService struct {
	metricsRepository db.MetricsDBRepository
}

func NewMetricsService(metricsRepository db.MetricsDBRepository) *MetricsService {
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
			s.metricsRepository.Add(ctx, metrics.MType, metrics.Name, *metrics.Value)
		} else {
			s.UpdateByTypeAndName(ctx, *metrics.Value, metrics.MType, metrics.Name)
		}
	case model.Counter:
		if existMetric == nil {
			s.metricsRepository.Add(ctx, metrics.MType, metrics.Name, float64(*metrics.Delta))
		} else {
			newValue := float64(*metrics.Delta + *existMetric.Delta)

			s.UpdateByTypeAndName(ctx, newValue, metrics.MType, metrics.Name)
		}
	}

	return nil
}

func (s *MetricsService) UpdateByTypeAndName(ctx context.Context, newValue float64, mType, name string) {
	s.metricsRepository.UpdateByTypeAndName(ctx, newValue, mType, name)
}
