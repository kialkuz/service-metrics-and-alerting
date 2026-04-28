package server

import (
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/file"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"time"
)

type MetricsFileService interface {
	SetMetric(name string, metric model.Metrics)
	Save()
}

var metricsList map[string]model.Metrics
var now time.Time

type FileService struct {
	repository    file.MetricsFileRepository
	storeInterval int
}

func NewFileService(repository file.MetricsFileRepository, storeInterval int) *FileService {
	now = time.Now()
	metricsList = make(map[string]model.Metrics)

	return &FileService{repository: repository, storeInterval: storeInterval}
}

func (p *FileService) SetMetric(name string, metric model.Metrics) {
	metricsList[name] = metric
}

func (p *FileService) SaveWithInterval() {
	for {
		if p.storeInterval != 0 {
			time.Sleep(time.Duration(p.storeInterval) * time.Second)
		}

		p.Save()
	}
}

func (p *FileService) Save() {
	p.repository.Truncate()

	for _, metric := range metricsList {
		p.repository.WriteMetric(&metric)
	}
}
