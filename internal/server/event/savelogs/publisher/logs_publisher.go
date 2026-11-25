package publisher

import (
	subscriberAudit "kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber/audit"
)

type LogsPublisher struct {
	observers map[string]subscriberAudit.AuditSaver
}

func NewLogsPublisher() LogsPublisher {
	return LogsPublisher{
		observers: make(map[string]subscriberAudit.AuditSaver),
	}
}

func (p *LogsPublisher) Register(o subscriberAudit.AuditSaver) {
	p.observers[o.GetAuditType()] = o
}

func (p *LogsPublisher) Shutdown() {
	for _, observer := range p.GetObservers() {
		p.Deregister(observer)
	}
}

func (p *LogsPublisher) Deregister(o subscriberAudit.AuditSaver) {
	o.Close()
	delete(p.observers, o.GetAuditType())
}

func (p *LogsPublisher) GetObservers() map[string]subscriberAudit.AuditSaver {
	return p.observers
}

func (p *LogsPublisher) GetObserver(auditType string) subscriberAudit.AuditSaver {
	observer, ok := p.observers[auditType]
	if ok {
		return observer
	}

	return nil
}
