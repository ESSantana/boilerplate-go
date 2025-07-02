package middlewares

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_count",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_error_count",
			Help: "Total number of HTTP errors",
		},
		[]string{"method", "path", "status"},
	)
)

func PrometheusInit() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ErrorCount)
}

func PrometheusMetricsHandler() fiber.Handler {
	return adaptor.HTTPHandler(promhttp.Handler())
}

func TrackMetrics() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		path := ctx.Path()
		method := ctx.Method()
		ctx.Next()
		status := ctx.Response().StatusCode()
		RequestCount.WithLabelValues(method, path, http.StatusText(status)).Inc()
		if status >= 400 {
			ErrorCount.WithLabelValues(method, path, http.StatusText(status)).Inc()
		}
		return nil
	}
}
