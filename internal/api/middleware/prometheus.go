package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "financeapp_http_requests_total",
			Help: "The total number of HTTP requests processed by the application.",
		},
		[]string{"route", "method", "status_code"},
	)

	ErrorCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "financeapp_http_requests_errors_total",
			Help: "The total number of HTTP ERROR response processed by the application.",
		},
		[]string{"route", "method", "status_code"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "financeapp_http_request_duration_seconds",
			Help:    "The delay time in seconds for each HTTP request.",
			Buckets: prometheus.DefBuckets, 
		},
		[]string{"route", "method"},
	)

	InFlightRequests = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "finance_http_in_flight_requests",
			Help: "The number of HTTP requests being processed in real-time.",
		},
	)
)

func PrometheusMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		InFlightRequests.Inc()
	
		defer InFlightRequests.Dec()

		start := time.Now()

		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		duration := time.Since(start)
		statusCode := ww.Status()
		routePattern := chi.RouteContext(r.Context()).RoutePattern()

		if routePattern == "" {
			routePattern = "not_found"
		}

		RequestDuration.WithLabelValues(routePattern, r.Method).Observe(duration.Seconds())
		RequestCount.WithLabelValues(routePattern, r.Method, strconv.Itoa(statusCode)).Inc()

		if statusCode >= 400 {
			ErrorCount.WithLabelValues(routePattern, r.Method, strconv.Itoa(statusCode)).Inc()
		}
	})
}
