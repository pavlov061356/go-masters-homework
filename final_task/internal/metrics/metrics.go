package metrics

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{0.1, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path", "status"},
	)

	ReviewsSentimentDistribution = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "reviews_sentiment_distribution",
			Help:    "Reviews sentiment distribution",
			Buckets: []float64{1, 2, 3},
		},
		[]string{},
	)

	SentimenterQueueLength = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sentimenter_queue_length",
		Help: "Length of the Sentimenter queue",
	})
)

// PrometheusMiddleware - middleware для сбора метрик
func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		status := strconv.Itoa(ww.Status())
		duration := time.Since(start).Seconds()

		httpRequestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, r.URL.Path, status).Observe(duration)
	})
}
