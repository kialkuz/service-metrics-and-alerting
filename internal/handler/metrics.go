package handler

import (
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	metricsService service.MetricsWriter
}

func NewMetricsHandler(metricsService service.MetricsWriter) *MetricsHandler {
	return &MetricsHandler{
		metricsService: metricsService,
	}
}

func (h *MetricsHandler) AddHandler(c *gin.Context) {
	contentType := c.GetHeader("Content-Type")
	if contentType != "text/plain" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect header type"})
		return
	}

	metricType := c.Param("type")

	if !slices.Contains(model.MType, metricType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect metric type"})
		return
	}

	var metrics model.Metrics
	metrics.MType = metricType
	metrics.Name = c.Param("name")
	value, err := strconv.ParseFloat(c.Param("value"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
	}
	metrics.Value = &value

	h.metricsService.Save(metrics)

	c.Status(http.StatusOK)
}
