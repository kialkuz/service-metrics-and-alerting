package server

import (
	"container/list"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"sync"
)

type AuditLogsList struct {
	list *list.List
	m    sync.Mutex
}

func NewAuditLogsList() *AuditLogsList {
	return &AuditLogsList{list: list.New()}
}

func (l *AuditLogsList) Add(log dto.MetricAudit) {
	l.m.Lock()
	defer l.m.Unlock()

	l.list.PushBack(log)
}

func (l *AuditLogsList) Pop() *dto.MetricAudit {
	l.m.Lock()
	defer l.m.Unlock()

	e := l.list.Front()
	if e == nil {
		return nil
	}

	logEvent := e.Value.(dto.MetricAudit)

	if l.list.Len() > 1 {
		l.list.MoveToFront(e.Next())
	}

	l.list.Remove(e)

	return &logEvent
}
