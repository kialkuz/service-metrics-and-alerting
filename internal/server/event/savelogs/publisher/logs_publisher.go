package publisher

import (
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"
)

type LogsPublisher struct {
	observers map[string]subscriber.AuditSaver
}

func NewLogsPublisher() LogsPublisher {
	return LogsPublisher{
		observers: make(map[string]subscriber.AuditSaver),
	}
}

func (p *LogsPublisher) Register(o subscriber.AuditSaver) {
	if p.observers == nil {
		p.observers = make((map[string]subscriber.AuditSaver))
	}
	p.observers[o.GetAuditType()] = o
}

func (p *LogsPublisher) Shutdown() {
	for _, observer := range p.observers {
		p.Deregister(observer)
	}
}

func (p *LogsPublisher) Deregister(o subscriber.AuditSaver) {
	o.Close()
	delete(p.observers, o.GetAuditType())
}

func (p *LogsPublisher) GetObservers() map[string]subscriber.AuditSaver {
	return p.observers
}
