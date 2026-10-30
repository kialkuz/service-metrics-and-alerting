package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"os"
	"sync"
)

type AuditFile struct {
	Audit
	m        sync.Mutex
	fileDesc *os.File
}

func NewAuditFile(savePath string) (*AuditFile, error) {
	if !service.IsFileName(savePath) {
		return nil, pkgErrors.ErrNotFIleName
	}

	audit := Audit{AuditType: savelogs.AuditFile}

	fileDesc, err := os.OpenFile(savePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &AuditFile{Audit: audit, fileDesc: fileDesc}, nil
}

func (a *AuditFile) Save(ctx context.Context, event savelogs.Event) error {
	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}
	jsonData = append(jsonData, '\n')

	a.m.Lock()
	defer a.m.Unlock()

	_, err = a.fileDesc.Write(jsonData)
	return err
}

func (a *AuditFile) Close() error {
	return a.fileDesc.Close()
}
