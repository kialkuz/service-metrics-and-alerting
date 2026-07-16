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

type FileService struct {
	mrw           sync.RWMutex
	repository    MetricsFileRepository
	storeInterval int
}

func NewFileService(repository MetricsFileRepository, storeInterval int) *FileService {
	now = time.Now()
	metricsList = make(map[string]model.Metrics)

	return &FileService{repository: repository, storeInterval: storeInterval}
}

func (p *FileService) SetMetric(name string, metric model.Metrics) {
	p.mrw.Lock()
	defer p.mrw.Unlock()

	metricsList[name] = metric
}

func (p *FileService) SaveWithInterval(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(p.storeInterval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case <-ctx.Done():
			return
		default:
			err := p.Save()
			if err != nil {
				log.Println("Error saving metrics:", err)
			}
		}
	}
}

func (p *FileService) Save() error {
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
