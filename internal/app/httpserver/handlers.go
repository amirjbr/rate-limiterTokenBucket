package httpserver

import (
	"TokenBucketRateLimiter/internal/core/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Handler struct {
	LimiterSvc *service.Service
}

func NewHandler(limitService *service.Service) *Handler {
	return &Handler{LimiterSvc: limitService}
}

func (h *Handler) LimitHandler(c *gin.Context) {
	userID := "01988728-1e86-750a-8379-c37c2a6e2285"
	ip := c.ClientIP()
	destinationService := "payment"
	method := "post"

	res := h.LimiterSvc.Limit(c, userID, ip, destinationService, method)
	if !res {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"Message": "rate limit exceeded",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Message": "allowed",
	})

}
