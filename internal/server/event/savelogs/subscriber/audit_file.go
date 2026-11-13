package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"os"
	"sync"
)

const AuditTypeFile = "file"

type AuditFile struct {
	Audit
	m    sync.Mutex
	file io.WriteCloser
}

func NewAuditFile(savePath string) (AuditSaver, error) {
	if !service.IsFileName(savePath) {
		return nil, pkgErrors.ErrNotFIleName
	}

	file, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &AuditFile{Audit: Audit{AuditType: AuditTypeFile}, file: file}, nil
}

func (a *AuditFile) Save(ctx context.Context, events <-chan savelogs.Event) {
	for event := range events {
		jsonData, err := json.Marshal(event)
		if err != nil {
			a.logger.Error(fmt.Errorf("ошибка сериализации JSON: %w, ip: %s", err, event.IPAddress))
			return
		}
		jsonData = append(jsonData, '\n')

		a.m.Lock()
		_, err = a.file.Write(jsonData)
		a.m.Unlock()

		if err != nil {
			a.logger.Error(fmt.Errorf("ошибка записи в файл логов аудита: %w, ip: %s", err, event.IPAddress))
			return
		}
	}
}

func (a *AuditFile) Close() error {
	return a.file.Close()
}
