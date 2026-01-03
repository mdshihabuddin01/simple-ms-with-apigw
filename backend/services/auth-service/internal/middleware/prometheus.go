package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	authHttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "auth_http_requests_total",
			Help: "Total number of HTTP requests for auth service",
		},
		[]string{"method", "endpoint", "status"},
	)

	authHttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_http_request_duration_seconds",
			Help:    "Auth service HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	authHttpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "auth_http_requests_in_flight",
			Help: "Number of auth service HTTP requests currently in flight",
		},
		[]string{"endpoint"},
	)

	authTokensIssuedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "auth_tokens_issued_total",
			Help: "Total number of auth tokens issued",
		},
	)

	authDatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "auth_database_query_duration_seconds",
			Help:    "Auth service database query duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"operation"},
	)
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		endpoint := c.FullPath()
		if endpoint == "" {
			endpoint = "unknown"
		}

		authHttpRequestsInFlight.WithLabelValues(endpoint).Inc()
		defer authHttpRequestsInFlight.WithLabelValues(endpoint).Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		authHttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		authHttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

func RecordAuthTokenIssued() {
	authTokensIssuedTotal.Inc()
}

func RecordDatabaseQuery(operation string, duration time.Duration) {
	authDatabaseQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}
