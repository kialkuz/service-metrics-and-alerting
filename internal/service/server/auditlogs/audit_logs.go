package auditlogs

import "kialkuz/service-metrics-and-alerting/internal/dto"

type AuditLogs interface {
	Add(log dto.MetricAudit)
	Pop() *dto.MetricAudit
}
