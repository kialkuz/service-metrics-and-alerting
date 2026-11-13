package auditlogs

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
)

type AuditLogs interface {
	Add(log dto.AuditLog)
	Pop(ctx context.Context, queryType string) (*dto.AuditLog, bool)
}
