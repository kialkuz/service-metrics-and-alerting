package contracts

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/model"
)

type MetricsService interface {
	SaveMetric(ctx context.Context, metrics model.Metrics) error
	SaveMetricList(ctx context.Context, metrics []model.Metrics) error
	Get(ctx context.Context, metricType, name string) (*model.Metrics, error)
	GetList(ctx context.Context) ([]model.Metrics, error)
	UpdateByTypeAndName(ctx context.Context, newValue float64, mType, name string) error
	GetGroupedByTypeAndName(ctx context.Context) (map[string]map[string]model.Metrics, error)
}
