package agent

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	"log"
	"time"
)

type MetricsAgent struct {
	metricsService service.MetricsAgentService
}

func NewMetricsAgent(metricsService service.MetricsAgentService) *MetricsAgent {
	return &MetricsAgent{
		metricsService: metricsService,
	}
}

func (a *MetricsAgent) Collect(reportInterval, pollInterval int) {
	now := time.Now()

	for {
		time.Sleep(time.Duration(pollInterval) * time.Second)
		if time.Now().After(now.Add(time.Duration(reportInterval) * time.Second)) {
			for fieldName, fieldValue := range a.metricsService.CollectCounter() {
				response, err := a.metricsService.Send(dto.Metrics{
					ID:    fieldName,
					MType: model.Counter,
					Delta: &fieldValue,
				})
				if err != nil {
					log.Println(err)
					continue
				}

				fmt.Println(response.Status)
			}

			for fieldName, fieldValue := range a.metricsService.CollectGauge() {
				response, err := a.metricsService.Send(dto.Metrics{
					ID:    fieldName,
					MType: model.Gauge,
					Value: &fieldValue,
				})
				if err != nil {
					log.Println(err)
					continue
				}

				fmt.Println(response.Status)
			}

			now = time.Now()
		}
	}
}
