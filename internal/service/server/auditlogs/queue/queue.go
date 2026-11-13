package queue

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
)

type Queue struct {
	queue chan *dto.AuditLog
}

func NewQuery(maxThreads int) Queue {
	return Queue{
		queue: make(chan *dto.AuditLog, maxThreads),
	}
}

func (q *Queue) Add(log dto.AuditLog) {
	q.queue <- &log
}

func (q *Queue) Pop(ctx context.Context) (*dto.AuditLog, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case v, ok := <-q.queue:
		return v, ok
	}
}
