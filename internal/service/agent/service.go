package service

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
)

type MetricsAgentService interface {
	Collect() map[string]map[string]float64
	Send(metricType string, name string, value float64) (resp *http.Response, err error)
}

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsService struct {
	url string
}

func NewMetricsService(url string) *MetricsService {
	return &MetricsService{url: url}
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
	response, err := http.Post(s.url+query, "text/plain", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}
