package subscriber

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
)

type AuditSaver interface {
	Save(ctx context.Context, event savelogs.Event) error
	GetAuditType() string
	Close() error
}

type Audit struct {
	AuditType string
}

func (a Audit) GetAuditType() string {
	return a.AuditType
}
