package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	subscriberAudit "kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber/audit"
	"kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs"
	"maps"
	"slices"
	"time"

	"go.uber.org/zap"
)

type Publisher interface {
	Register(subscriberAudit.AuditSaver)
	Shutdown()
	Deregister(subscriberAudit.AuditSaver)
	GetObservers() map[string]subscriberAudit.AuditSaver
	Save(ctx context.Context)
}

type LogsPublisherQuery struct {
	LogsPublisher
	auditLogger auditlogs.AuditLogger
	logger      *zap.SugaredLogger
	eventsChans map[string]chan savelogs.Event
}

func NewLogsPublisherQuery(auditLogger auditlogs.AuditLogger, logger *zap.SugaredLogger) Publisher {
	return &LogsPublisherQuery{
		LogsPublisher: NewLogsPublisher(),
		auditLogger:   auditLogger,
		logger:        logger,
	}
}

func (p *LogsPublisherQuery) Register(o subscriberAudit.AuditSaver) {
	o.SetLogger(p.logger)
	p.LogsPublisher.Register(o)
}

func (p *LogsPublisherQuery) Shutdown() {
	for _, o := range p.GetObservers() {
		close(p.eventsChans[o.GetAuditType()])
	}

	p.LogsPublisher.Shutdown()
}

func (p *LogsPublisherQuery) Deregister(o subscriberAudit.AuditSaver) {
	eventsChan := p.eventsChans[o.GetAuditType()]

	if len(eventsChan) > 0 {
		for {
			if len(eventsChan) > 0 {
				time.Sleep(time.Duration(1 * time.Second))
				continue
			}

			break
		}
	}

	p.LogsPublisher.Deregister(o)
}

func (p *LogsPublisherQuery) Save(ctx context.Context) {
	p.eventsChans = p.getChans()

	for {
		select {
		case <-ctx.Done():
			for _, observer := range p.GetObservers() {
				p.Deregister(observer)
			}
			return
		default:
			for name, observer := range p.GetObservers() {
				go observer.Save(ctx, p.eventsChans[name])
			}

			p.sendToSubscriber(ctx, p.eventsChans)
		}
	}
}

func (p *LogsPublisherQuery) getChans() map[string]chan savelogs.Event {
	chans := make(map[string]chan savelogs.Event, len(p.GetObservers()))

	for name := range p.observers {
		chans[name] = make(chan savelogs.Event, 10)
	}

	return chans
}

func (p *LogsPublisherQuery) sendToSubscriber(ctx context.Context, eventsChans map[string]chan savelogs.Event) {
	queryList := slices.Collect(maps.Keys(p.GetObservers()))

	for {
		for _, queryType := range queryList {
			log, ok := p.auditLogger.Pop(ctx, queryType)
			if !ok {
				return
			}

			eventsChans[queryType] <- savelogs.Event{
				TS:        time.Now().Unix(),
				Metrics:   log.ListName,
				IPAddress: log.IP,
			}
		}
	}
}
