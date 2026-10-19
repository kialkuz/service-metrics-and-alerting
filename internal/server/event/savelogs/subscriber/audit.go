package subscriber

import "kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"

type AuditSaver interface {
	Save(logsData savelogs.Event) error
	GetAuditType() string
}

type Audit struct {
	AuditType string
	SavePath  string
}

func (a Audit) GetAuditType() string {
	return a.AuditType
}
