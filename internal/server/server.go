package server

import (
	"kialkuz/service-metrics-and-alerting/internal/handler"
	"kialkuz/service-metrics-and-alerting/internal/router"
	"net/http"
	"time"
)

func NewServer(handler *handler.MetricsHandler, serverPort string) *http.Server {
	return &http.Server{
		Addr:         ":" + serverPort,
		Handler:      router.Init(handler),
		ReadTimeout:  time.Duration(5 * time.Second),
		WriteTimeout: time.Duration(10 * time.Second),
		IdleTimeout:  time.Duration(15 * time.Second),
	}
}
