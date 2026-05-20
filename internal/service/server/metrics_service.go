package server

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

type MetricsRepository interface {
	Add(ctx context.Context, typeValue, name string, value float64) error
	AddList(ctx context.Context, metrics []model.Metrics) error
	SaveList(ctx context.Context, metricsForInsert []model.Metrics, metricsForUpdate []model.Metrics) error
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

func (s *MetricsService) GetGroupedByTypeAndName(ctx context.Context) (map[string]map[string]model.Metrics, error) {
	items, err := s.metricsRepository.GetList(ctx)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]map[string]model.Metrics)

	for _, item := range items {
		if _, exists := metrics[item.MType]; !exists {
			metrics[item.MType] = make(map[string]model.Metrics)
		}

		metrics[item.MType][item.Name] = item
	}

	return metrics, nil
}

func (s *MetricsService) SaveMetric(ctx context.Context, metrics model.Metrics) error {
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

func (s *MetricsService) SaveMetricList(ctx context.Context, metrics []model.Metrics) error {
	existMetrics, _ := s.GetGroupedByTypeAndName(ctx)

	var metricsForInsert []model.Metrics
	var metricsForUpdate []model.Metrics

	for _, metric := range metrics {
		existMetric, exists := existMetrics[metric.MType][metric.Name]

		if !exists {
			metricsForInsert = append(metricsForInsert, metric)
		} else {
			if metric.MType == model.Counter {
				newValue := *metric.Delta + *existMetric.Delta
				existMetric.Delta = &newValue
			} else {
				if existMetric.Value == metric.Value {
					continue
				}

				existMetric.Value = metric.Value
			}

			metricsForUpdate = append(metricsForUpdate, existMetric)
		}
	}

	return s.metricsRepository.SaveList(ctx, metricsForInsert, metricsForUpdate)
}

func (s *MetricsService) UpdateByTypeAndName(ctx context.Context, newValue float64, mType, name string) error {
	return s.metricsRepository.UpdateByTypeAndName(ctx, newValue, mType, name)
}
