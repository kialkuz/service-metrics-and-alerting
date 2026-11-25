package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/logger"
)

type Queue struct {
	queue  chan *dto.AuditLog
	logger logger.Logger
}

func NewQuery(maxThreads int, logger logger.Logger) Queue {
	return Queue{
		queue:  make(chan *dto.AuditLog, maxThreads),
		logger: logger,
	}
}

func (q *Queue) Add(log dto.AuditLog) {
	select {
	case q.queue <- &log:
	default:
		b, err := json.Marshal(log.ListName)
		if err != nil {
			q.logger.Error(fmt.Sprintf("очередь заполнена, ошибка сериализации ListName: %v", err))
			return
		}
		q.logger.Warn("очередь заполнена, не записан лог, ListName: " + string(b))
	}
}

func (q *Queue) Pop(ctx context.Context) (*dto.AuditLog, bool) {
	select {
	case <-ctx.Done():
		return nil, false
	case v, ok := <-q.queue:
		return v, ok
	}
}
