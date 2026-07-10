package agent

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service/compress"
	signService "kialkuz/service-metrics-and-alerting/internal/service/sign"
	"log"
	"maps"
	"math/rand/v2"
	"net/http"
	"reflect"
	"runtime"
)

type MetricsAgentService interface {
	CollectCounter() map[string]int64
	CollectGauge() map[string]float64
	SendSingleMetric(value dto.Metrics) (*http.Response, error)
	SendListMetrics(value []dto.Metrics) (*http.Response, error)
	SendJSONBody(jsonData []byte, path string) (*http.Response, error)
	SendPostQuery(metricType, name, value string) (*http.Response, error)
}

//go:generate go run go.uber.org/mock/mockgen -source=service.go -destination=mocks/service_mock.go -package=mocks -typed
type MetricsService struct {
	url string
	key string
}

func NewMetricsService(url, key string) *MetricsService {
	return &MetricsService{url: url, key: key}
}

func (s *MetricsService) CollectCounter() map[string]int64 {
	fields := map[string]int64{
		"PollCount": 1,
	}

	return fields
}

func (s *MetricsService) CollectGauge() map[string]float64 {
	fields := map[string]float64{
		"RandomValue": rand.Float64(),
	}
	maps.Insert(fields, maps.All(s.collectMemStats()))

	return fields
}

func (s *MetricsService) collectMemStats() map[string]float64 {
	memStats := runtime.MemStats{}
	runtime.ReadMemStats(&memStats)

	r := reflect.ValueOf(memStats)

	statsFields := make(map[string]float64, len(model.StatsFields))
	for _, field := range model.StatsFields {
		f := r.FieldByName(field)
		if !f.IsValid() {
			log.Printf("Поле %s не валидное", field)
			continue
		}

		switch f.Kind() {
		case reflect.Uint32, reflect.Uint64:
			statsFields[field] = float64(f.Uint())
		case reflect.Float64:
			statsFields[field] = f.Float()
		}
	}

	return statsFields
}

func (s *MetricsService) SendSingleMetric(value dto.Metrics) (*http.Response, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	response, err := s.SendJSONBody(jsonData, "/update/")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *MetricsService) SendListMetrics(value []dto.Metrics) (*http.Response, error) {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("error marshalling JSON: %s", err.Error())
	}

	response, err := s.SendJSONBody(jsonData, "/updates/")
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *MetricsService) SendJSONBody(jsonData []byte, path string) (*http.Response, error) {
	b, err := compress.MakeGzip(jsonData)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, s.url+path, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Add("Content-Encoding", "gzip")
	if s.key != "" {
		request.Header.Add("HashSHA256", hex.EncodeToString(signService.Generate(jsonData, s.key)))
	}

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}

func (s *MetricsService) SendPostQuery(metricType, name, value string) (*http.Response, error) {
	query := fmt.Sprintf("/update/%s/%s/%s", metricType, name, value)
	response, err := http.Post(s.url+query, "text/plain", nil)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	return response, nil
}
