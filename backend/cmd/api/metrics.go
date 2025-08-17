package main

import (
	"net/http"
	"runtime"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Declare the metrics we want to expose.
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"}, // Labels
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method", "path"}, // Labels
	)

	appInfo = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "app_info",
		Help: "Information about the application.",
		ConstLabels: prometheus.Labels{
			"go_version": runtime.Version(),
		},
	})
)

// metricsHandler creates a handler to expose the Prometheus metrics.
func (app *application) metricsHandler() http.Handler {
	// Register the standard Go process and build info collectors.
	prometheus.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	prometheus.MustRegister(prometheus.NewGoCollector())

	// Set the app_info gauge to 1 to indicate the app is running.
	appInfo.Set(1)

	return promhttp.Handler()
}
