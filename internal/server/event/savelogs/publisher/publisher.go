package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"time"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const maxThreads = 10

type Publisher interface {
	Register(subscriber.AuditSaver)
	Deregister(subscriber.AuditSaver)
	SaveWithInterval(ctx context.Context)
	Save(ctx context.Context, log dto.MetricAudit)
}

type LogsPublisher struct {
	observers     map[string]subscriber.AuditSaver
	auditLogsList *service.AuditLogsList
	logger        *zap.SugaredLogger
	semaphore     *Semaphore
}

func NewLogsPublisher(auditLogsList *service.AuditLogsList, logger *zap.SugaredLogger) Publisher {
	return &LogsPublisher{
		observers:     make(map[string]subscriber.AuditSaver),
		auditLogsList: auditLogsList,
		logger:        logger,
		semaphore:     NewSemaphore(maxThreads),
	}
}

func (p *LogsPublisher) Register(o subscriber.AuditSaver) {
	if p.observers == nil {
		p.observers = make((map[string]subscriber.AuditSaver))
	}
	p.observers[o.GetAuditType()] = o
}

func (p *LogsPublisher) Deregister(o subscriber.AuditSaver) {
	o.Close()
	delete(p.observers, o.GetAuditType())
}

// Если событий будет очень много, то отправляем в clickhouse
// clickhouse работает сразу с тысячами сообщений, поэтому накапливаем их в очереди лога
// и потом сохраняем в бд. Для этого передосать saveLogs и дописать функцию для работы с бд.
func (p *LogsPublisher) SaveWithInterval(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			for _, observer := range p.observers {
				p.Deregister(observer)
			}
			return
		case <-ticker.C:
			log := p.auditLogsList.Pop()
			if log != nil {
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
}

// Пока используем запись логов на рантайме,
func (p *LogsPublisher) Save(ctx context.Context, log dto.MetricAudit) {
	err := p.saveLogs(ctx, savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   log.ListName,
		IPAddress: log.IP,
	})
	if err != nil {
		p.logger.Error("Error saving metrics logs:", err)
	}
}

func (p *LogsPublisher) saveLogs(ctx context.Context, event savelogs.Event) error {
	p.semaphore.Acquire()
	defer p.semaphore.Release()

	g, ctx := errgroup.WithContext(ctx)

	for _, observer := range p.observers {
		g.Go(func() error {
			return observer.Save(ctx, event)
		})
	}

	return g.Wait()
}
