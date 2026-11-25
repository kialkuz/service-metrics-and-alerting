package handler

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"net/http"
	"time"

	exampleConfig "kialkuz/service-metrics-and-alerting/example"
)

func MetricsInit() error {
	ctx := context.Background()
	path := "update/" + model.Counter + "/" + "TestCounter" + "/100"
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		exampleConfig.ServerURL+"/"+path,
		nil,
	)
	if err != nil {
		return err
	}

	client := &http.Client{
		Timeout: 2 * time.Second, // общий дедлайн на весь запрос
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
