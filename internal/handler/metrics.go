package handler

import (
	"context"
	"errors"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/service"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const timeout = 10

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
<<<<<<< iter3
		log.Println("11111111111111111")
=======
>>>>>>> v2
		log.Println(err)
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	newValue := c.Param("value")
	if newValue == "" {
<<<<<<< iter3
		log.Println("2222222222222222222222")
=======
>>>>>>> v2
		log.Println(err)
		c.JSON(http.StatusBadRequest, errors.New("empty metric value"))
		return
	}

	var metrics model.Metrics
	metrics.MType = metricType
	metrics.Name = name
	value, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
<<<<<<< iter3
		log.Println("3333333333333333333333333")
=======
>>>>>>> v2
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
		return
	}
	metrics.Value = &value

<<<<<<< iter3
	err = h.metricsService.Save(metrics)
	if err != nil {
		log.Println("4444444444444444444444")
=======
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	err = h.metricsService.Save(ctx, metrics)
	if err != nil {
>>>>>>> v2
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

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	metric, err := h.metricsService.Get(ctx, metricType, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, *metric.Value)
}

func (h *MetricsHandler) GetListHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	items, err := h.metricsService.GetList(ctx)
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
