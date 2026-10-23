package publisher

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"log"
	"time"

	"golang.org/x/sync/errgroup"
)

type Publisher interface {
	Register(subscriber.AuditSaver)
	Deregister(subscriber.AuditSaver)
	SaveWithInterval(ctx context.Context)
}

type LogsPublisher struct {
	observers    map[string]subscriber.AuditSaver
	auditService *service.AuditService
}

func NewLogsPublisher(auditService *service.AuditService) Publisher {
	return &LogsPublisher{
		observers:    make(map[string]subscriber.AuditSaver),
		auditService: auditService,
	}
}

func (p *LogsPublisher) Register(o subscriber.AuditSaver) {
	if p.observers == nil {
		p.observers = make((map[string]subscriber.AuditSaver))
	}
	p.observers[o.GetAuditType()] = o
}

func (p *LogsPublisher) Deregister(o subscriber.AuditSaver) {
	delete(p.observers, o.GetAuditType())
}

func (p *LogsPublisher) SaveWithInterval(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(5) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			err := p.save(savelogs.Event{
				TS:        time.Now().Unix(),
				Metrics:   p.auditService.GetMetrics(),
				IPAddress: p.auditService.GetIP(),
			})
			if err != nil {
				log.Println("Error saving metrics logs:", err)
			}
		}
	}
}

func (p *LogsPublisher) save(event savelogs.Event) error {
	g := new(errgroup.Group)

	for _, observer := range p.observers {
		g.Go(func() error {
			return observer.Save(event)
		})
	}

	if err := g.Wait(); err != nil {
		return err
	}

	return nil
}
