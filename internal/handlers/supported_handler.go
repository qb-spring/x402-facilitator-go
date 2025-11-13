package handlers

import (
	"net/http"
	"x402-facilitator-go/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SupportedHandler handles supported schemes and networks requests
type SupportedHandler struct {
	supportedService *service.SupportedService
	logger           *zap.Logger
}

// NewSupportedHandler creates a new SupportedHandler
func NewSupportedHandler(supportedService *service.SupportedService, logger *zap.Logger) *SupportedHandler {
	return &SupportedHandler{
		supportedService: supportedService,
		logger:           logger,
	}
}

// Supported handles GET /supported requests
func (h *SupportedHandler) Supported(c *gin.Context) {
	response := h.supportedService.Supported()
	c.JSON(http.StatusOK, response)
}
