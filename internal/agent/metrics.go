package agent

import (
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/service"
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
		metrics := a.metricsService.Collect()

		time.Sleep(time.Duration(pollInterval) * time.Second)
		if time.Now().After(now.Add(time.Duration(reportInterval) * time.Second)) {
			for metricType, metricList := range metrics {
				for name, value := range metricList {
					response, err := a.metricsService.Send(metricType, name, value)
					if err != nil {
						log.Println(err)
						return
					}

					fmt.Println(response.Status)
				}
			}

			now = time.Now()
		}
	}
}
