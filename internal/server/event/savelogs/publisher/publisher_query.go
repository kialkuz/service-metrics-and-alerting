package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"
	"kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs"
	"maps"
	"slices"
	"time"

	"go.uber.org/zap"
)

const maxHandlers = 3

var events map[string]chan savelogs.Event

type Publisher interface {
	Register(subscriber.AuditSaver)
	Shutdown()
	Deregister(subscriber.AuditSaver)
	GetObservers() map[string]subscriber.AuditSaver
	Save(ctx context.Context)
}

type LogsPublisherQuery struct {
	LogsPublisher
	auditLogs auditlogs.AuditLogs
	logger    *zap.SugaredLogger
}

func NewLogsPublisherQuery(auditLogs auditlogs.AuditLogs, logger *zap.SugaredLogger) Publisher {
	return &LogsPublisherQuery{
		LogsPublisher: NewLogsPublisher(),
		auditLogs:     auditLogs,
		logger:        logger,
	}
}

func (p *LogsPublisherQuery) Register(o subscriber.AuditSaver) {
	o.SetLogger(p.logger)
	p.LogsPublisher.Register(o)
}

func (p *LogsPublisherQuery) Save(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, observer := range p.GetObservers() {
				p.Deregister(observer)
			}
			return
		default:
			events = p.GetChans(maxHandlers)

			for w := 1; w <= maxHandlers; w++ {
				for name, observer := range p.GetObservers() {
					go observer.Save(ctx, events[name])
				}
			}

			p.sendToSubscriber(ctx)
		}
	}
}

func (p *LogsPublisherQuery) sendToSubscriber(ctx context.Context) {
	queryList := slices.Collect(maps.Keys(p.GetObservers()))

	for {
		for _, queryType := range queryList {
			log, ok := p.auditLogs.Pop(ctx, queryType)
			if !ok {
				return
			}

			events[queryType] <- savelogs.Event{
				TS:        time.Now().Unix(),
				Metrics:   log.ListName,
				IPAddress: log.IP,
			}
		}
	}
}
