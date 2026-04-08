package service

import (
	"context"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/config"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
)

type MetricsServerService interface {
	Save(ctx context.Context, name model.Metrics) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
}

type MetricsAgentService interface {
	Collect() map[string]map[string]float64
	Send(metricType string, name string, value float64) (resp *http.Response, err error)
}

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsService struct {
	metricsRepository db.MetricsRepository
}

func NewMetricsService(metricsRepository db.MetricsRepository) *MetricsService {
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
	switch metrics.MType {
	case model.Gauge:
		if err := s.updateGauge(ctx, metrics); err != nil {
			return err
		}
	case model.Counter:
		if err := s.updateCounter(ctx, metrics); err != nil {
			return err
		}
	}

	return nil
}

func (s *MetricsService) updateGauge(ctx context.Context, metrics model.Metrics) error {
	existMetric, err := s.metricsRepository.Get(ctx, metrics.MType, metrics.Name)
	if err != nil {
		if err := s.metricsRepository.Add(ctx, metrics); err != nil {
			return err
		}
	} else {
		if err := s.metricsRepository.Update(ctx, *metrics.Value, existMetric.ID); err != nil {
			return fmt.Errorf("repo SaveMetric: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) updateCounter(ctx context.Context, metrics model.Metrics) error {
	existMetric, err := s.metricsRepository.Get(ctx, metrics.MType, metrics.Name)
	if err != nil {
		if err := s.metricsRepository.Add(ctx, metrics); err != nil {
			return err
		}
	} else {
		newValue := *metrics.Value + (*existMetric.Value)

		if err := s.metricsRepository.Update(ctx, newValue, existMetric.ID); err != nil {
			return fmt.Errorf("repo SaveMetric: %w", err)
		}
	}

	return nil
}

func (s *MetricsService) Collect() map[string]map[string]float64 {
	metrics := make(map[string]map[string]float64)
	metrics[model.Counter] = make(map[string]float64)
	metrics[model.Counter]["RandomValue"] = rand.Float64()
	metrics[model.Gauge] = make(map[string]float64)
	metrics[model.Gauge] = s.collectMemStats()

	return metrics
}

func (s *MetricsService) collectMemStats() map[string]float64 {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	r := reflect.ValueOf(memStats)

	statsFields := make(map[string]float64)
	for _, field := range model.StatsFields {
		switch reflect.Indirect(r).FieldByName(field).Type().Name() {
		case "uint32":
		case "uint64":
			statsFields[field] = float64(reflect.Indirect(r).FieldByName(field).Uint())
		case "float64":
			statsFields[field] = float64(reflect.Indirect(r).FieldByName(field).Float())
		}
	}

	return statsFields
}

func (s *MetricsService) Send(metricType string, name string, value float64) (*http.Response, error) {
	query := fmt.Sprintf("/update/%s/%s/%f", metricType, name, value)
	response, err := http.Post(config.GetURL()+query, "text/plain", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}
