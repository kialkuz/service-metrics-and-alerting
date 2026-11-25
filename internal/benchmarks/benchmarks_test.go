package benchmarks

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/memory"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	"testing"
)

func BenchmarkStorageOperations(b *testing.B) {
	ctx := context.Background()

	delta := int64(1)
	metrics := []model.Metrics{{
		Name:  "TestCounter",
		MType: model.Counter,
		Delta: &delta,
	}}

	metric := metrics[0]

	b.ResetTimer()

	metricsService := service.NewMetricsService(memory.NewMemoryStorage())

	for b.Loop() {
		err := metricsService.SaveMetricList(ctx, metrics)
		if err != nil {
			b.Fatal(err)
			return
		}

		_, err = metricsService.Get(ctx, metric.MType, metric.Name)
		if err != nil {
			b.Fatal(err)
			return
		}
	}
}
