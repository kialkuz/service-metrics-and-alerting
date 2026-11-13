package subscriber

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"kialkuz/service-metrics-and-alerting/internal/service/request"
	"kialkuz/service-metrics-and-alerting/internal/service/retrier"
	"kialkuz/service-metrics-and-alerting/pkg/errorclassifier"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"net/http"
)

const (
	AuditTypeURL  = "url"
	retryAttempts = 3
)

type AuditURL struct {
	Audit
	savePath string
	retrier  retrier.Retrier
}

func NewAuditURL(savePath string) (AuditSaver, error) {
	if !service.IsURL(savePath) {
		return nil, pkgErrors.ErrNotURL
	}

	return &AuditURL{
		Audit:    Audit{AuditType: AuditTypeURL},
		savePath: savePath,
		retrier: retrier.NewRetrierService(
			retrier.RetryConfig{Attempts: retryAttempts},
			errorclassifier.NewNetworkErrorClassifier(),
		),
	}, nil
}

func (a AuditURL) Save(ctx context.Context, events <-chan savelogs.Event) {
	for event := range events {
		jsonData, err := json.MarshalIndent(event, "", " ")
		if err != nil {
			a.logger.Error(fmt.Errorf("ошибка сериализации JSON: %w", err))
			return
		}

		err = a.sendJSONBody(ctx, jsonData)
		if err != nil {
			a.logger.Error(err)
			return
		}
	}
}

func (a AuditURL) sendJSONBody(ctx context.Context, jsonData []byte) error {
	client := &http.Client{}

	return a.retrier.WithRetry(ctx, func(ctx context.Context) error {
		requestData, err := http.NewRequestWithContext(ctx, http.MethodPost, a.savePath, bytes.NewReader(jsonData))
		if err != nil {
			return err
		}

		requestData.Header.Set("Content-Type", "application/json")

		return request.Send(client, requestData)
	})
}

func (a AuditURL) Close() error {
	return nil
}
