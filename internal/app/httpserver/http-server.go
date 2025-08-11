package httpserver

import (
	"TokenBucketRateLimiter/internal/app/httpserver/middleware"
	"TokenBucketRateLimiter/internal/core/service"
	"TokenBucketRateLimiter/pkg/observeutil"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Engine            *gin.Engine
	Handlers          *Handler
	LimiterSvc        *service.Service
	metricsMiddleware *middleware.MetricsMiddleware
}

func NewHttpServer(handler *Handler, metrics *observeutil.ProMetrics) *HttpServer {
	var httpServer HttpServer
	httpServer.metricsMiddleware = middleware.NewMetricsMiddleware(metrics)
	httpServer.createEngine()
	httpServer.Handlers = handler
	httpServer.baseRoutes()
	httpServer.Engine.Use(httpServer.metricsMiddleware.Handle)
	httpServer.initialRoutes()

	return &httpServer
}

func (h *HttpServer) createEngine() {
	h.Engine = gin.Default()
}

func (h *HttpServer) baseRoutes() {
	h.Engine.GET("/metrics", gin.WrapH(observeutil.PrometheusHandler()))
}

func (h *HttpServer) initialRoutes() {
	h.Engine.GET("/rate-limiter/limit", h.Handlers.LimitHandler)
}
