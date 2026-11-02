package subscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/storage"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"net/http"
)

const retryAttempts = 3

type AuditURL struct {
	Audit
	retry      storage.RetryConfig
	classifier errorclassifier.ErrorClassifier
	savePath   string
}

func NewAuditURL(savePath string) (*AuditURL, error) {
	if !service.IsURL(savePath) {
		return nil, pkgErrors.ErrNotURL
	}

	audit := Audit{AuditType: savelogs.AuditURL}

	return &AuditURL{
		Audit:      audit,
		retry:      storage.RetryConfig{Attempts: retryAttempts},
		classifier: errorclassifier.NewNetworkErrorClassifier(),
		savePath:   savePath,
	}, nil
}

func (a AuditURL) Save(ctx context.Context, event savelogs.Event) error {
	jsonData, err := json.MarshalIndent(event, "", " ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	return a.sendJSONBody(ctx, jsonData)
}

func (a AuditURL) sendJSONBody(ctx context.Context, jsonData []byte) error {
	client := &http.Client{}
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, a.savePath, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	_, err = storage.WithRetry(ctx, a.retry, a.classifier, func(ctx context.Context) (struct{}, error) {
		response, err := client.Do(request)
		if err != nil {
			return struct{}{}, err
		}
		defer response.Body.Close()
		return struct{}{}, err
	})

	return err
}

func (a AuditURL) Close() error {
	return nil
}
