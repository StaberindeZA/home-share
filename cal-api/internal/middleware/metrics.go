package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type HTTPMetrics struct {
	RequestsTotal   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
}

func NewHTTPMetrics(reg prometheus.Registerer) *HTTPMetrics {
	factory := promauto.With(reg)
	return &HTTPMetrics{
		RequestsTotal: factory.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests processed.",
			},
			[]string{"path", "status"},
		),
		RequestDuration: factory.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Duration of HTTP requests in seconds.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"path"},
		),
	}
}

func (m HTTPMetrics) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// A simple custom response writer wrapper to capture the status code
		wrappedWriter := &ResponseWriterDelegator{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(wrappedWriter, r)

		duration := time.Since(start).Seconds()
		statusStr := fmt.Sprintf("%d", wrappedWriter.statusCode)

		// r.Pattern keeps the route clean (e.g., "/v1/home/{slug}/mates")
		// Fallback to r.URL.Path if no pattern is matched
		pathPattern := r.Pattern
		if pathPattern == "" {
			pathPattern = r.URL.Path
		}

		// Record the metrics
		m.RequestsTotal.WithLabelValues(pathPattern, statusStr).Inc()
		m.RequestDuration.WithLabelValues(pathPattern).Observe(duration)
	})
}
