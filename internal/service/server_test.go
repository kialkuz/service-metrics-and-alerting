package service

import (
	"math/rand"
	"testing"

	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db/mocks"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/pkg/errors"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddMetricSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	metricName := "test_name"
	metricValue := rand.Float64()

	mockRepo.EXPECT().GetMetric(metricName).Return(nil, errors.ErrNotFound)
	mockRepo.EXPECT().AddMetric(gomock.Any()).Return(nil)

	modelMetrics := model.Metrics{
		MType: model.Counter,
		Name:  metricName,
		Value: &metricValue,
	}
	err := service.Save(modelMetrics)

	assert.NoError(t, err)
}

func TestUpdateMetricSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	currentMetricValue := rand.Float64()

	metricName := "test_name"
	newMetricValue := rand.Float64()

	modelMetrics := model.Metrics{
		ID:    "1",
		MType: model.Counter,
		Name:  metricName,
		Value: &newMetricValue,
	}

	mockRepo.EXPECT().GetMetric(metricName).Return(&model.Metrics{
		ID:    modelMetrics.ID,
		MType: modelMetrics.MType,
		Name:  modelMetrics.Name,
		Value: &currentMetricValue,
	}, nil)
	mockRepo.EXPECT().UpdateMetric(gomock.Any(), gomock.Any()).Return(nil)

	err := service.Save(modelMetrics)

	assert.NoError(t, err)
}
