package queue

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"

	"go.uber.org/zap"
)

const maxThreads = 10

type AuditLogsQueueList struct {
	queues map[string]Queue
}

func NewAuditLogsQueueList(queuesName []string, logger *zap.SugaredLogger) *AuditLogsQueueList {
	queues := make(map[string]Queue, len(queuesName))

	for _, queueName := range queuesName {
		queues[queueName] = NewQuery(maxThreads, logger)
	}

	return &AuditLogsQueueList{
		queues: queues,
	}
}

func (l *AuditLogsQueueList) Add(log dto.AuditLog) {
	for queueType := range l.queues {
		queue := l.queues[queueType]
		queue.Add(log)
	}
}

func (l *AuditLogsQueueList) Pop(ctx context.Context, queueType string) (*dto.AuditLog, bool) {
	queue := l.queues[queueType]

	return queue.Pop(ctx)
}
