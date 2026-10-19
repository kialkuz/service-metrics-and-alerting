package subscriber

import (
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/service"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"os"
)

type AuditFile struct {
	Audit
}

func NewAuditFile(savePath string) (*AuditFile, error) {
	if !service.IsFileName(savePath) {
		return nil, pkgErrors.ErrNotFIleName
	}

	audit := Audit{AuditType: savelogs.AuditFile, SavePath: savePath}

	return &AuditFile{Audit: audit}, nil
}

func (a AuditFile) Save(event savelogs.Event) error {
	jsonData, err := json.MarshalIndent(event, "", " ")
	if err != nil {
		return fmt.Errorf("ошибка сериализации JSON: %w", err)
	}

	f, err := os.OpenFile(
		a.Audit.SavePath,
		os.O_WRONLY|os.O_CREATE|os.O_APPEND,
		0644,
	)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(jsonData)
	return err
}
