package server

import (
	"context"

	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service/server/mocks"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	ctx := context.TODO()

	metrics := model.Metrics{
		ID:    1,
		MType: model.Counter,
		Name:  "PollCount",
		Delta: new(int64),
	}
	*metrics.Delta = 100
	mockRepo.EXPECT().Get(ctx, gomock.Any(), gomock.Any()).Return(&metrics, nil)

	result, err := service.Get(ctx, model.Counter, "PollCount")

	assert.NoError(t, err)
	assert.Equal(t, result, &metrics)
}

func TestGetNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	ctx := context.TODO()

	mockRepo.EXPECT().Get(ctx, gomock.Any(), gomock.Any()).Return(nil, pkgErrors.ErrNotFound)

	result, err := service.Get(ctx, model.Counter, "PollCount")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestGetListSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockMetricsRepository(ctrl)
	service := NewMetricsService(mockRepo)

	ctx := context.TODO()

	metrics := []model.Metrics{
		{
			ID:    1,
			MType: model.Counter,
			Name:  "PollCount",
			Delta: new(int64),
		},
		{
			ID:    2,
			MType: model.Gauge,
			Name:  "HeapIdle",
			Value: new(float64),
		},
	}
	*metrics[0].Delta = 100
	*metrics[1].Value = 200.5
	mockRepo.EXPECT().GetList(ctx).Return(metrics, nil)

	result, err := service.GetList(ctx)

	require.NoError(t, err)
	require.Len(t, result, 2)
	require.Equal(t, metrics, result)
}
