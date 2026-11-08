package auditlogs

import (
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/service"
)

const maxThreads = 10

type AuditLogsQuery struct {
	query     chan *dto.MetricAudit
	semaphore *service.Semaphore
}

func NewAuditLogsQuery() *AuditLogsQuery {
	return &AuditLogsQuery{
		query:     make(chan *dto.MetricAudit, 10),
		semaphore: service.NewSemaphore(maxThreads),
	}
}

func (l *AuditLogsQuery) Add(log dto.MetricAudit) {
	l.semaphore.Acquire()
	defer l.semaphore.Release()

	select {
	case l.query <- &log:
	default:
	}
}

func (l *AuditLogsQuery) Pop() *dto.MetricAudit {
	return <-l.query
}
