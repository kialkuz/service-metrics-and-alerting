package handler

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	service "kialkuz/service-metrics-and-alerting/internal/service/server"
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
	switch metricType {
	case model.Counter:
		value, err := strconv.ParseInt(newValue, 10, 64)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
			return
		}
		metrics.Delta = &value
	case model.Gauge:
		value, err := strconv.ParseFloat(newValue, 64)
		if err != nil {
			log.Println(err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
			return
		}
		metrics.Value = &value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	err = h.metricsService.Save(ctx, metrics)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *MetricsHandler) UpdateHandler(c *gin.Context) {
	var request dto.Metrics

	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := request.ID
	metricType := request.MType

	httpCode, err := h.check(metricType, id)
	if err != nil {
		log.Println(err)
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	var metrics model.Metrics
	metrics.MType = metricType
	metrics.Name = id

	switch metricType {
	case model.Counter:
		if request.Delta == nil {
			error := fmt.Errorf("с типом %s нужно передавать параметр delta", model.Counter)
			log.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}

		metrics.Delta = request.Delta
	case model.Gauge:
		if request.Value == nil {
			error := fmt.Errorf("с типом %s нужно передавать параметр value", model.Gauge)
			log.Println(error)
			c.JSON(http.StatusBadRequest, gin.H{"error": error.Error()})
			return
		}

		metrics.Value = request.Value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()

	err = h.metricsService.Save(ctx, metrics)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *MetricsHandler) GetMetricHandler(c *gin.Context) {
	var request dto.Metrics

	if err := c.BindJSON(&request); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный JSON"})
		return
	}

	httpCode, err := h.check(request.MType, request.ID)
	if err != nil {
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	metric, err := h.metricsService.Get(ctx, request.MType, request.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var responseDto dto.Metrics
	responseDto.ID = metric.Name
	responseDto.MType = metric.MType

	if metric.MType == model.Counter {
		responseDto.Delta = metric.Delta
	} else {
		responseDto.Value = metric.Value
	}

	c.JSON(http.StatusOK, responseDto)
}

func (h *MetricsHandler) GetMetricValueHandler(c *gin.Context) {
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

	if metric.MType == model.Counter {
		c.JSON(http.StatusOK, *metric.Delta)
	} else {
		c.JSON(http.StatusOK, *metric.Value)
	}
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
		if item.Delta != nil {
			metrics = append(metrics, dto.MetricView{
				Name:  item.Name,
				Delta: *item.Delta,
			})
		} else {
			metrics = append(metrics, dto.MetricView{
				Name:  item.Name,
				Value: *item.Value,
			})
		}
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
