package subscriber

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/logger"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
)

//go:generate go run go.uber.org/mock/mockgen -source=audit.go -destination=mocks/audit_mock.go -package=mocks -typed
type AuditSaver interface {
	GetAuditType() string
	Close() error
	Save(ctx context.Context, events <-chan savelogs.Event)
	SetLogger(logger logger.Logger)
}

type Audit struct {
	AuditType string
	logger    logger.Logger
}

func (a Audit) GetAuditType() string {
	return a.AuditType
}

func (a *Audit) SetLogger(logger logger.Logger) {
	a.logger = logger
}
