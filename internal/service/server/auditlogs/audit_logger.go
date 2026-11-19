package auditlogs

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
)

type AuditLogger interface {
	Add(log dto.AuditLog)
	Pop(ctx context.Context, queryType string) (*dto.AuditLog, bool)
}
