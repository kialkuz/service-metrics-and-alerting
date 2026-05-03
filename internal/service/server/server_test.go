package server

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db/mocks"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/pkg/errors"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestAddMetricSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsDBRepository(ctrl)
	service := NewMetricsService(mockRepo)

	metricType := model.Counter
	metricName := "test_name"
	metricValue := rand.Int63()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mockRepo.EXPECT().Get(ctx, metricType, metricName).Return(nil, errors.ErrNotFound)
	mockRepo.EXPECT().Add(ctx, gomock.Any(), gomock.Any(), gomock.Any())

	modelMetrics := model.Metrics{
		MType: model.Counter,
		Name:  metricName,
		Delta: &metricValue,
	}
	err := service.Save(ctx, modelMetrics)

	assert.NoError(t, err)
}

func TestUpdateMetricSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsDBRepository(ctrl)
	service := NewMetricsService(mockRepo)

	currentMetricValue := rand.Int63()

	metricType := model.Counter
	metricName := "test_name"
	newMetricValue := rand.Int63()

	modelMetrics := model.Metrics{
		ID:    1,
		MType: model.Counter,
		Name:  metricName,
		Delta: &newMetricValue,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	mockRepo.EXPECT().Get(ctx, metricType, metricName).Return(&model.Metrics{
		ID:    modelMetrics.ID,
		MType: modelMetrics.MType,
		Name:  modelMetrics.Name,
		Delta: &currentMetricValue,
	}, nil)
	mockRepo.EXPECT().UpdateByTypeAndName(ctx, gomock.Any(), gomock.Any(), gomock.Any())

	err := service.Save(ctx, modelMetrics)

	assert.NoError(t, err)
}
