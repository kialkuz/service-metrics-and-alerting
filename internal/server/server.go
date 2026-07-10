package init

import (
	"kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/router"
	"kialkuz/service-metrics-and-alerting/internal/server/handler"
	"net/http"
	"time"
)

func NewServer(handler *handler.MetricsHandler, config *server.Config) *http.Server {
	return &http.Server{
		Addr:         ":" + config.ServerPort,
		Handler:      router.Init(handler, config),
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}
