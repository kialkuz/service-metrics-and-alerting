package publisher

import (
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
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

func (p *LogsPublisher) GetObserver(auditType string) subscriber.AuditSaver {
	observer, ok := p.observers[auditType]
	if ok {
		return observer
	}

	return nil
}

func (p *LogsPublisher) GetChans(maxThreads int) map[string]chan savelogs.Event {
	chans := make(map[string]chan savelogs.Event, len(p.observers))

	for name := range p.observers {
		chans[name] = make(chan savelogs.Event, maxThreads)
	}

	return chans
}
