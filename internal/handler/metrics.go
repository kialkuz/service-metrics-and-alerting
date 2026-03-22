package handler

import (
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"net/http"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MetricsHandler struct {
	metricsService service.MetricsServerService
}

func NewMetricsHandler(metricsService service.MetricsServerService) *MetricsHandler {
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
	name := c.Param("name")

	err := h.check(metricType, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var metrics model.Metrics
	metrics.MType = metricType
	metrics.Name = name
	value, err := strconv.ParseFloat(c.Param("value"), 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
	}
	metrics.Value = &value

	h.metricsService.Save(metrics)

	c.Status(http.StatusOK)
}

func (h *MetricsHandler) GetMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")

	err := h.check(metricType, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	metric, err := h.metricsService.Get(metricType, name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, metric)
}

func (h *MetricsHandler) GetListHandler(c *gin.Context) {
	items, err := h.metricsService.GetList()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var metrics []dto.MetricView
	for _, item := range items {
		metrics = append(metrics, dto.MetricView{
			Name:  item.Name,
			Value: *item.Value,
		})
	}

	c.HTML(http.StatusOK, "metrics_list.html", gin.H{"metrics": metrics})
}

func (h *MetricsHandler) check(metricType, name string) error {
	if !slices.Contains(model.MType, metricType) {
		return errors.New("Incorrect metric type")
	}

	if name == "" {
		return errors.New("Empty metric name")
	}

	return nil
}
