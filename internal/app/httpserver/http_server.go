package httpserver

import (
	"TokenBucketRateLimiter/internal/core/service"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Engine     *gin.Engine
	Handlers   *Handler
	LimiterSvc *service.Service
}

func NewHttpServer(handler *Handler) *HttpServer {
	var httpServer HttpServer

	httpServer.createEngine()
	httpServer.Handlers = handler
	httpServer.initialRoutes()

	return &httpServer
}

func (h *HttpServer) createEngine() {
	h.Engine = gin.Default()
}

func (h *HttpServer) initialRoutes() {
	routes := h.Engine.Group("/rateLimiter")

	routes.GET("/limit", h.Handlers.LimitHandler)
}
