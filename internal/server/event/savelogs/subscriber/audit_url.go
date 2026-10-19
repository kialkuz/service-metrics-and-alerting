package subscriber

import (
	"bytes"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"net/http"
)

type AuditUrl struct {
	Audit
}

func NewAuditUrl(savePath string) (*AuditUrl, error) {
	if !service.IsFileName(savePath) {
		return nil, pkgErrors.ErrNotFIleName
	}

	audit := Audit{AuditType: savelogs.AuditUrl, SavePath: savePath}

	return &AuditUrl{Audit: audit}, nil
}

func (a AuditUrl) Save(event savelogs.Event) error {
	jsonData, err := json.MarshalIndent(event, "", " ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	return a.sendJSONBody(jsonData, a.SavePath)
}

func (a AuditUrl) sendJSONBody(jsonData []byte, path string) error {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()

	return nil
}
