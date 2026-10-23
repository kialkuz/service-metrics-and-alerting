package server

import (
	"sync"
)

type AuditService struct {
	mrw         sync.Mutex
	metricsList map[string]bool
	ip          string
}

func NewAuditService() *AuditService {
	return &AuditService{metricsList: make(map[string]bool)}
}

func (a *AuditService) SetMetric(name string) {
	a.mrw.Lock()
	defer a.mrw.Unlock()

	a.metricsList[name] = true
}

func (a *AuditService) SetIP(ip string) {
	a.ip = ip
}

func (a *AuditService) GetIP() string {
	return a.ip
}

func (a *AuditService) GetMetrics() []string {
	var metrics []string

	for metric := range a.metricsList {
		metrics = append(metrics, metric)
	}

	return metrics
}
