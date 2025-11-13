package handlers

import (
	"net/http"
	"x402-facilitator-go/internal/middleware"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// SettleHandler handles settlement requests
type SettleHandler struct {
	settleService *service.SettleService
	logger        *zap.Logger
}

// NewSettleHandler creates a new SettleHandler
func NewSettleHandler(settleService *service.SettleService, logger *zap.Logger) *SettleHandler {
	return &SettleHandler{
		settleService: settleService,
		logger:        logger,
	}
}

// Settle handles POST /settle requests
func (h *SettleHandler) Settle(c *gin.Context) {
	requestLogger := middleware.GetRequestLogger(c, h.logger)

	var request models.SettleRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		requestLogger.Warn("Invalid request body",
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Call the settlement service
	ctx := c.Request.Context()
	response := h.settleService.Settle(ctx, &request)

	c.JSON(http.StatusOK, response)
}
