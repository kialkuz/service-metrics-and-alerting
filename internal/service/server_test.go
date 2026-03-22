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

	metricType := model.Counter
	metricName := "test_name"
	metricValue := rand.Float64()

	mockRepo.EXPECT().Get(metricType, metricName).Return(nil, errors.ErrNotFound)
	mockRepo.EXPECT().Add(gomock.Any()).Return(nil)

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

	metricType := model.Counter
	metricName := "test_name"
	newMetricValue := rand.Float64()

	modelMetrics := model.Metrics{
		ID:    "1",
		MType: model.Counter,
		Name:  metricName,
		Value: &newMetricValue,
	}

	mockRepo.EXPECT().Get(metricType, metricName).Return(&model.Metrics{
		ID:    modelMetrics.ID,
		MType: modelMetrics.MType,
		Name:  modelMetrics.Name,
		Value: &currentMetricValue,
	}, nil)
	mockRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)

	err := service.Save(modelMetrics)

	assert.NoError(t, err)
}
