package service

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"strconv"
)

type MetricsRepository interface {
	AddMetric(metrics model.Metrics) error
	UpdateMetric(value float64, id int) error
	GetMetric(name string) (*model.Metrics, error)
}

type MetricsService struct {
	metricsRepository MetricsRepository
}

func NewMetricsService(metricsRepository MetricsRepository) *MetricsService {
	return &MetricsService{metricsRepository: metricsRepository}
}

func (s *MetricsService) Save(metrics model.Metrics) error {
	switch metrics.MType {
	case model.Gauge:
		if err := s.metricsRepository.AddMetric(metrics); err != nil {
			return err
		}
	case model.Counter:
		if err := s.updateMetric(metrics); err != nil {
			return err
		}
	}

	return nil
}

func (s *MetricsService) updateMetric(metrics model.Metrics) error {
	existMetric, err := s.metricsRepository.GetMetric(metrics.Name)
	if err != nil {
		return err
	}
	id, err := strconv.Atoi(existMetric.ID)
	if err != nil {
		return err
	}
	newValue := *metrics.Value + (*existMetric.Value)

	if err := s.metricsRepository.UpdateMetric(newValue, id); err != nil {
		return fmt.Errorf("repo SaveMetric: %w", err)
	}

	return nil
}
