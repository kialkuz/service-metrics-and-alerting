package server

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"log"
	"sync"
	"time"
)

type MetricsFileRepository interface {
	CreateTemp() error
	WriteMetric(metric *model.Metrics) error
	SaveOriginal() error
}

var metricsList map[string]model.Metrics
var now time.Time

type LogsSaverService struct {
	mrw           sync.RWMutex
	repository    MetricsFileRepository
	storeInterval int
}

func NewLogsSaverService(repository MetricsFileRepository, storeInterval int) *LogsSaverService {
	now = time.Now()
	metricsList = make(map[string]model.Metrics)

	return &LogsSaverService{repository: repository, storeInterval: storeInterval}
}

func (p *LogsSaverService) SetMetric(name string, metric model.Metrics) {
	p.mrw.Lock()
	defer p.mrw.Unlock()

	metricsList[name] = metric
}

func (p *LogsSaverService) SaveWithInterval(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.storeInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			err := p.save()
			if err != nil {
				log.Println("Error saving metrics:", err)
			}
		}
	}
}

func (p *LogsSaverService) save() error {
	p.mrw.RLock()
	defer p.mrw.RUnlock()

	if len(metricsList) != 0 {
		err := p.repository.CreateTemp()
		if err != nil {
			return err
		}

		for _, metric := range metricsList {
			err = p.repository.WriteMetric(&metric)
			if err != nil {
				return err
			}
		}

		err = p.repository.SaveOriginal()
		if err != nil {
			return err
		}
	}

	return nil
}
