package main

import (
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// appMetrics holds all our custom prometheus metrics.
type appMetrics struct {
	requestsTotal    *prometheus.CounterVec
	requestsDuration *prometheus.HistogramVec
	appInfo          prometheus.Gauge
}

// newMetrics creates, registers, and returns our metrics struct.
func newMetrics(reg *prometheus.Registry) appMetrics {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "path", "status"},
	)
	requestsDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Duration of HTTP requests.",
		},
		[]string{"method", "path"},
	)
	appInfo := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "app_info",
		Help: "Information about the application.",
		ConstLabels: prometheus.Labels{
			"go_version": runtime.Version(),
		},
	})

	// Register all metrics with the provided custom registry.
	reg.MustRegister(requestsTotal, requestsDuration, appInfo)
	reg.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	reg.MustRegister(prometheus.NewGoCollector())

	appInfo.Set(1)

	return appMetrics{
		requestsTotal:    requestsTotal,
		requestsDuration: requestsDuration,
		appInfo:          appInfo,
	}
}

// withMetrics is the middleware that records metrics for each request.
func (app *application) withMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		// Use the metrics from the app.metrics struct, not global variables.
		app.metrics.requestsTotal.WithLabelValues(r.Method, r.URL.Path, strconv.Itoa(rw.statusCode)).Inc()
		app.metrics.requestsDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
	})
}
