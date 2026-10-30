package agent

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/agent/collectalgo"
	service "kialkuz/service-metrics-and-alerting/internal/service/agent"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"time"
)

const defaultCountAdditionalReplaySendMetrics = 3

type MetricsAgent struct {
	metricsService service.MetricsAgentService
	metricsConfig  MetricsConfig
}

func NewMetricsAgent(
	metricsService service.MetricsAgentService,
	metricsConfig MetricsConfig,
) *MetricsAgent {
	return &MetricsAgent{
		metricsService: metricsService,
		metricsConfig:  metricsConfig,
	}
}

func (a *MetricsAgent) CollectAndSend(ctx context.Context) error {
	now := time.Now()

	ticker := time.NewTicker(time.Duration(a.metricsConfig.PollInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return pkgErrors.ErrAppStopped
		case <-ticker.C:
			if time.Now().After(now.Add(time.Duration(a.metricsConfig.ReportInterval) * time.Second)) {
				collectalgo, err := a.getCollectAlgo()
				if err != nil {
					return err
				}

				collectalgo.CollectAndSend(ctx)

				now = time.Now()
			}
		}
	}
}

func (a *MetricsAgent) getCollectConfig() collectalgo.CollectConfig {
	rateLimit := a.metricsConfig.RateLimit

	return collectalgo.CollectConfig{
		CountAdditionalReplaySendMetrics: defaultCountAdditionalReplaySendMetrics,
		AdditionalReplaySendInterval: map[int]int{
			1: 1,
			2: 3,
			3: 5,
		},
		RateLimit: rateLimit,
	}
}

func (a *MetricsAgent) getCollectAlgo() (collectalgo.CollectAlgo, error) {
	collectConfig := a.getCollectConfig()

	switch a.metricsConfig.CollectAlgo {
	case collectalgo.CollectSingle:
		return collectalgo.NewSingle(a.metricsService, collectConfig), nil
	case collectalgo.CollectAll:
		return collectalgo.NewAll(a.metricsService, collectConfig), nil
	default:
		return nil, errors.New("unknown collect algo")
	}
}
