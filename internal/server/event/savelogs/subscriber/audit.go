package subscriber

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
)

type AuditSaver interface {
	GetAuditType() string
	Close() error
	Save(ctx context.Context, event savelogs.Event) error
}

type Audit struct {
	AuditType string
}

func (a Audit) GetAuditType() string {
	return a.AuditType
}
