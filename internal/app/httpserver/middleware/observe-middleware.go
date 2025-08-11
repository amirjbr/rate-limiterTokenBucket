package middleware

import (
	"TokenBucketRateLimiter/pkg/observeutil"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type MetricsMiddleware struct {
	metrics *observeutil.ProMetrics
}

func NewMetricsMiddleware(metrics *observeutil.ProMetrics) *MetricsMiddleware {
	return &MetricsMiddleware{
		metrics: metrics,
	}
}

func (m *MetricsMiddleware) Handle(ctx *gin.Context) {
	start := time.Now()

	m.metrics.HTTP.IncInFlight()
	ctx.Next()
	m.metrics.HTTP.DecInFlight()

	durationMs := float64(time.Since(start).Milliseconds())

	handler := ctx.FullPath()
	if handler == "" {
		handler = ctx.Request.URL.Path
	}
	method := ctx.Request.Method
	statusCode := ctx.Writer.Status()
	statusStr := strconv.Itoa(statusCode)

	m.metrics.HTTP.IncRequest(handler, method, statusStr)
	m.metrics.HTTP.ObserveDuration(handler, method, durationMs)

	switch {
	case statusCode >= 200 && statusCode < 300:
		m.metrics.HTTP.Inc2xx(handler, method)
	case statusCode >= 400 && statusCode < 500:
		m.metrics.HTTP.Inc4xx(handler, method)
	case statusCode >= 500:
		m.metrics.HTTP.Inc5xx(handler, method)
	}

}
