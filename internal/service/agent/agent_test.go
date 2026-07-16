package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCollectMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMetricsService("test_url", "")

	gaugeMetrics, err := service.CollectGauge()
	assert.NotEmpty(t, service.CollectCounter())
	assert.Empty(t, err)
	assert.NotEmpty(t, gaugeMetrics)
}
