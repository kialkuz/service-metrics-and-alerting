package handler

import (
	"context"
	"errors"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/dto"
	"kialkuz/service-metrics-and-alerting/internal/model"
	pkgContracts "kialkuz/service-metrics-and-alerting/pkg/contracts"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
	"net/http"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const timeout = 10

type MetricsFileService interface {
	SetMetric(name string, metric model.Metrics)
	Save() error
}

type AuditService interface {
	SetIP(ip string)
	SetMetric(name string)
}

type Pinger interface {
	Ping(ctx context.Context) error
}

type MetricsHandler struct {
	metricsService     pkgContracts.MetricsService
	metricsFileService MetricsFileService
	metricsDBService   Pinger
	auditService       AuditService
}

func NewMetricsHandler(
	metricsService pkgContracts.MetricsService,
	metricsFileService MetricsFileService,
	metricsDBService Pinger,
	auditService AuditService,
) *MetricsHandler {
	return &MetricsHandler{
		metricsService:     metricsService,
		metricsFileService: metricsFileService,
		metricsDBService:   metricsDBService,
		auditService:       auditService,
	}
}

// PingDB godoc
// @Summary Ping db
// @Description Запрос доступа бд
// @Success 200
// @Failure 500
// @Router /ping [get]
func (h *MetricsHandler) PingDB(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Second)
	defer cancel()
	if err := h.metricsDBService.Ping(ctx); err != nil {
		if errors.Is(err, pkgErrors.ErrNotInitDB) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Not initialized db"})
			return
		}

		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{})
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

// AddHandler godoc
// @Summary Запрос добавления метрики
// @Description Запрос добавления метрики
// @Param type path string true "Тип метрики" default(gauge)
// @Param name path string true "Имя метрики" default(HeapSys)
// @Param value path number true "Значение метрики" default(1000)
// @Success 200
// @Failure 400
// @Router /update/{type}/{name}/{value} [post]
func (h *MetricsHandler) AddHandler(c *gin.Context) {
	metricType := c.Param("type")
	name := c.Param("name")

	httpCode, err := h.check(metricType, name)
	if err != nil {
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	newValue := c.Param("value")
	if newValue == "" {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
			return
		}
		metrics.Delta = &value
	case model.Gauge:
		value, err := strconv.ParseFloat(newValue, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Incorrect value type"})
			return
		}
		metrics.Value = &value
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout*time.Second)
	defer cancel()
	err = h.metricsService.SaveMetric(ctx, metrics)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	c.Status(http.StatusOK)
}

// UpdateHandler godoc
// @Summary Запрос изменения метрики
// @Description Запрос изменения метрики
// @Param req body dto.Metrics true "Тело запроса (id, type, delta/value)"
// @Success 200
// @Failure 400
// @Router /update/ [post]
func (h *MetricsHandler) UpdateHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var requestMetric dto.Metrics

	if err := c.ShouldBindJSON(&requestMetric); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	httpCode, err := h.check(requestMetric.MType, requestMetric.ID)
	if err != nil {
		c.Error(err)
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	var metrics model.Metrics
	metrics.MType = requestMetric.MType
	metrics.Name = requestMetric.ID

	httpCode, err = h.validateUpdateRequest(requestMetric)
	if err != nil {
		c.Error(err)
		c.JSON(httpCode, gin.H{"error": err.Error()})
		return
	}

	switch requestMetric.MType {
	case model.Counter:
		metrics.Delta = requestMetric.Delta
	case model.Gauge:
		metrics.Value = requestMetric.Value
	}

	err = h.metricsService.SaveMetric(ctx, metrics)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	item, err := h.metricsService.Get(ctx, requestMetric.MType, requestMetric.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	h.metricsFileService.SetMetric(metrics.Name, *item)
	h.auditService.SetIP(c.RemoteIP())
	h.auditService.SetMetric(metrics.Name)

	c.Status(http.StatusOK)
}

// UpdatesHandler godoc
// @Summary Запрос изменения метрик
// @Description Запрос изменения метрик
// @Param req body []dto.Metrics true "Тело запроса (id, type, delta/value)"
// @Success 200
// @Failure 400
// @Failure 404
// @Failure 0
// @Router /updates/ [post]
func (h *MetricsHandler) UpdatesHandler(c *gin.Context) {
	ctx := c.Request.Context()

	var requestMetricsList []dto.Metrics

	if err := c.ShouldBindJSON(&requestMetricsList); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var metrics []model.Metrics

	for _, requestMetric := range requestMetricsList {
		httpCode, err := h.check(requestMetric.MType, requestMetric.ID)
		if err != nil {
			c.Error(err)
			c.JSON(httpCode, gin.H{"error": err.Error()})
			return
		}

		httpCode, err = h.validateUpdateRequest(requestMetric)
		if err != nil {
			c.Error(err)
			c.JSON(httpCode, gin.H{"error": err.Error()})
			return
		}

		metric := &model.Metrics{
			MType: requestMetric.MType,
			Name:  requestMetric.ID,
		}

		switch requestMetric.MType {
		case model.Counter:
			metric.Delta = requestMetric.Delta
		case model.Gauge:
			metric.Value = requestMetric.Value
		}

		metrics = append(metrics, *metric)
	}

	err := h.metricsService.SaveMetricList(ctx, metrics)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	dbMetrics, err := h.metricsService.GetList(ctx)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Metric not saved"})
		return
	}

	for _, dbMetric := range dbMetrics {
		h.metricsFileService.SetMetric(dbMetric.Name, dbMetric)
		h.auditService.SetIP(c.RemoteIP())
		h.auditService.SetMetric(dbMetric.Name)
	}

	c.Status(http.StatusOK)
}

func (h *MetricsHandler) validateUpdateRequest(metrics dto.Metrics) (httpCode int, err error) {
	switch metrics.MType {
	case model.Counter:
		if metrics.Delta == nil {
			return http.StatusBadRequest, fmt.Errorf("с типом %s нужно передавать параметр delta", model.Counter)
		}
	case model.Gauge:
		if metrics.Value == nil {
			return http.StatusBadRequest, fmt.Errorf("с типом %s нужно передавать параметр value", model.Gauge)
		}
	default:
		return 0, errors.New("unknown metric type")
	}

	return 0, nil
}

// GetMetricHandler godoc
// @Summary Запрос получения метрики
// @Description Запрос получения метрики
// @Param req body dto.GetMetric true "Тело запроса (id, type)"
// @Success 200
// @Failure 404
// @Failure 500
// @Failure 0
// @Router /value/ [post]
func (h *MetricsHandler) GetMetricHandler(c *gin.Context) {
	var request dto.GetMetric

	if err := c.ShouldBindJSON(&request); err != nil {
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
		c.Error(err)

		if metric == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "metric not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
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

// GetMetricValueHandler godoc
// @Summary Запрос получения значения метрики
// @Description Запрос получения значения метрики
// @Param type path string true "Тип метрики" Enums(counter, gauge)
// @Param name path string true "Имя метрики" example(PollCount)
// @Success 200
// @Failure 404
// @Failure 0
// @Router /value/{type}/{name} [get]
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

// GetListHandler godoc
// @Summary Запрос получения списка метрик
// @Description Запрос получения списка метрик
// @Success 200
// @Failure 404
// @Router / [get]
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
