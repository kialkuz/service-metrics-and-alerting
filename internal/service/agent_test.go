package service

import (
	"testing"

	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db/mocks"
	"kialkuz/service-metrics-and-alerting/internal/model"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCollectMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	metrics := service.Collect()

	assert.NotEmpty(t, metrics[model.Counter])
	assert.NotEmpty(t, metrics[model.Gauge])
}
