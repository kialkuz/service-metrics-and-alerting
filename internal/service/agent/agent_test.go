package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestCollectMetrics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	service := NewMetricsService("test_url")

	assert.NotEmpty(t, service.CollectCounter())
	assert.NotEmpty(t, service.CollectGauge())
}
