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

// initMetrics registers all our custom metrics with the application's registry.
func (app *application) initMetrics() {
	// Register application-specific metrics.
	app.metrics_reg.MustRegister(httpRequestsTotal)
	app.metrics_reg.MustRegister(httpRequestDuration)
	app.metrics_reg.MustRegister(appInfo)

	// Register the standard Go process and build info collectors.
	app.metrics_reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	app.metrics_reg.MustRegister(prometheus.NewGoCollector())

	// Set the initial value for the app_info gauge.
	appInfo.Set(1)
}

// metricsHandler creates a handler that serves metrics from our CUSTOM registry.
func (app *application) metricsHandler() http.Handler {
	return promhttp.HandlerFor(app.metrics_reg, promhttp.HandlerOpts{})
}
