package handler

import (
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
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
	metricType := c.Param("type")
	name := c.Param("name")

	httpCode, err := h.check(metricType, name)
	if err != nil {
		log.Println(err)
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	newValue := c.Param("value")
	if newValue == "" {
		log.Println(err)
		c.JSON(http.StatusBadRequest, errors.New("empty metric value"))
		return
	}

	var metrics model.Metrics
	metrics.MType = metricType
	metrics.Name = name
	value, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
		return
	}
	metrics.Value = &value

	err = h.metricsService.Save(metrics)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *MetricsHandler) GetMetricHandler(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")

	httpCode, err := h.check(metricType, name)
	if err != nil {
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	metric, err := h.metricsService.Get(metricType, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *metric.Value)
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

func (h *MetricsHandler) check(metricType, name string) (int, error) {
	if !slices.Contains(model.MType, metricType) {
		return http.StatusBadRequest, errors.New("incorrect metric type")
	}

	if name == "" {
		return http.StatusNotFound, errors.New("empty metric name")
	}

	return 0, nil
}
