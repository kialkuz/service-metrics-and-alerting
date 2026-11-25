package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	subscriberAudit "kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber/audit"
	"kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs"
	"maps"
	"slices"
	"sync"
	"time"

	"go.uber.org/zap"
)

type LogsPublisherQuery struct {
	LogsPublisher
	auditLogger auditlogs.AuditLogger
	logger      *zap.SugaredLogger
	wg          *sync.WaitGroup
}

func NewLogsPublisherQuery(auditLogger auditlogs.AuditLogger, logger *zap.SugaredLogger) *LogsPublisherQuery {
	return &LogsPublisherQuery{
		LogsPublisher: NewLogsPublisher(),
		auditLogger:   auditLogger,
		logger:        logger,
		wg:            &sync.WaitGroup{},
	}
}

func (p *LogsPublisherQuery) Register(o subscriberAudit.AuditSaver) {
	o.SetLogger(p.logger)
	p.LogsPublisher.Register(o)
}

func (p *LogsPublisherQuery) Shutdown(eventsChans map[string]chan savelogs.Event) {
	for _, ch := range eventsChans {
		close(ch)
	}

	p.wg.Wait()

	p.LogsPublisher.Shutdown()
}

func (p *LogsPublisherQuery) Save(ctx context.Context) {
	var eventsChans map[string]chan savelogs.Event

	for {
		select {
		case <-ctx.Done():
			p.Shutdown(eventsChans)
			return
		default:
			if eventsChans == nil {
				eventsChans = p.getChans()
				p.startSubscribers(ctx, eventsChans)
			}

			p.sendToSubscriber(ctx, eventsChans)
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

func (p *LogsPublisherQuery) startSubscribers(
	ctx context.Context,
	eventsChans map[string]chan savelogs.Event,
) {
	for name, observer := range p.GetObservers() {
		p.wg.Add(1)
		p.wg.Go(func() {
			defer p.wg.Done()
			observer.Save(ctx, eventsChans[name])
		})
	}
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
