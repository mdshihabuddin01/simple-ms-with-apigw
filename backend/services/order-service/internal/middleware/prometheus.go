package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	orderHttpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "order_http_requests_total",
			Help: "Total number of HTTP requests for order service",
		},
		[]string{"method", "endpoint", "status"},
	)

	orderHttpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_http_request_duration_seconds",
			Help:    "Order service HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	orderHttpRequestsInFlight = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "order_http_requests_in_flight",
			Help: "Number of order service HTTP requests currently in flight",
		},
		[]string{"endpoint"},
	)

	ordersCreatedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "orders_created_total",
			Help: "Total number of orders created",
		},
	)

	orderDatabaseQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "order_database_query_duration_seconds",
			Help:    "Order service database query duration in seconds",
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

		orderHttpRequestsInFlight.WithLabelValues(endpoint).Inc()
		defer orderHttpRequestsInFlight.WithLabelValues(endpoint).Dec()

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		orderHttpRequestsTotal.WithLabelValues(method, endpoint, status).Inc()
		orderHttpRequestDuration.WithLabelValues(method, endpoint).Observe(duration)
	}
}

func RecordOrderCreated() {
	ordersCreatedTotal.Inc()
}

func RecordDatabaseQuery(operation string, duration time.Duration) {
	orderDatabaseQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}
