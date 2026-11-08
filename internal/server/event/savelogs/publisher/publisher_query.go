package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"
	"kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

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
			log := p.auditLogs.Pop()
			err := p.saveLogs(ctx, savelogs.Event{
				TS:        time.Now().Unix(),
				Metrics:   log.ListName,
				IPAddress: log.IP,
			})
			if err != nil {
				p.logger.Error("Error saving metrics logs:", err)
			}
		}
	}
}

func (p *LogsPublisherQuery) saveLogs(ctx context.Context, event savelogs.Event) error {
	g, ctx := errgroup.WithContext(ctx)

	for _, observer := range p.GetObservers() {
		g.Go(func() error {
			return observer.Save(ctx, event)
		})
	}

	return g.Wait()
}
