package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	exampleConfig "kialkuz/service-metrics-and-alerting/example"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"net/http"
	"time"
)

func ExampleMetricsHandler_AddHandler() {
	ctx := context.Background()
	path := "update/" + model.Counter + "/" + "PollCount" + "/100"
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		exampleConfig.ServerURL+"/"+path,
		nil,
	)
	if err != nil {
		panic(fmt.Errorf("error create request: %w", err))
	}

	resp, err := makeRequest(req)
	if err != nil {
		panic(fmt.Errorf("error execute request: %w", err))
	}
	defer resp.Body.Close()

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func ExampleMetricsHandler_UpdateHandler() {
	metric := dto.Metrics{
		ID:    "TestCounter",
		MType: model.Counter,
		Delta: new(int64),
	}
	*metric.Delta = 100

	body, err := json.Marshal(metric)
	if err != nil {
		panic(fmt.Errorf("error encode request: %w", err))
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		exampleConfig.ServerURL+"/update/",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(fmt.Errorf("error create request: %w", err))
	}

	resp, err := makeRequest(req)
	if err != nil {
		panic(fmt.Errorf("error execute request: %w", err))
	}
	defer resp.Body.Close()

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func ExampleMetricsHandler_UpdatesHandler() {
	metric := []dto.Metrics{{
		ID:    "TestCounter",
		MType: model.Counter,
		Delta: new(int64),
	}}
	*metric[0].Delta = 100

	body, err := json.Marshal(metric)
	if err != nil {
		panic(fmt.Errorf("error encode request: %w", err))
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		exampleConfig.ServerURL+"/updates/",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(fmt.Errorf("error create request: %w", err))
	}

	resp, err := makeRequest(req)
	if err != nil {
		panic(fmt.Errorf("error execute request: %w", err))
	}
	defer resp.Body.Close()

	// fmt.Println(string(body))
	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func ExampleMetricsHandler_GetMetricHandler() {
	MetricsInit()

	metric := dto.GetMetric{
		ID:    "TestCounter",
		MType: model.Counter,
	}

	body, err := json.Marshal(metric)
	if err != nil {
		panic(fmt.Errorf("error encode request: %w", err))
	}

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		exampleConfig.ServerURL+"/value/",
		bytes.NewBuffer(body),
	)
	if err != nil {
		panic(fmt.Errorf("error create request: %w", err))
	}

	resp, err := makeRequest(req)
	if err != nil {
		panic(fmt.Errorf("error execute request: %w", err))
	}
	defer resp.Body.Close()

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func ExampleMetricsHandler_GetMetricValueHandler() {
	MetricsInit()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		exampleConfig.ServerURL+"/value/"+model.Counter+"/TestCounter",
		nil,
	)
	if err != nil {
		panic(fmt.Errorf("error create request: %w", err))
	}

	resp, err := makeRequest(req)
	if err != nil {
		panic(fmt.Errorf("error execute request: %w", err))
	}
	defer resp.Body.Close()

	fmt.Printf("Статус ответа: %d\n", resp.StatusCode)
	// Output: Статус ответа: 200
}

func makeRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{
		Timeout: 2 * time.Second, // общий дедлайн на весь запрос
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp, fmt.Errorf("error execute request: %w", err)
	}

	return resp, nil
}
