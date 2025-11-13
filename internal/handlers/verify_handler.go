package handlers

import (
	"net/http"

	"x402-facilitator-go/internal/middleware"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/service"
	"x402-facilitator-go/pkg/errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// VerifyHandler handles verification requests
type VerifyHandler struct {
	verifyService *service.VerifyService
	logger        *zap.Logger
}

// NewVerifyHandler creates a new VerifyHandler
func NewVerifyHandler(verifyService *service.VerifyService, logger *zap.Logger) *VerifyHandler {
	return &VerifyHandler{
		verifyService: verifyService,
		logger:        logger,
	}
}

// Verify handles POST /verify requests
func (h *VerifyHandler) Verify(c *gin.Context) {
	requestLogger := middleware.GetRequestLogger(c, h.logger)

	var request models.VerifyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		requestLogger.Warn("Invalid request body", zap.Error(err))
		// Extract payer from request if possible, otherwise empty string
		c.JSON(http.StatusBadRequest, models.VerifyResponse{
			IsValid:       false,
			InvalidReason: errors.ErrorInvalidPayload.Code(),
			Payer:         request.PaymentPayload.Payload.Authorization.From,
		})
		return
	}

	// Call the verification service with context
	ctx := c.Request.Context()
	response := h.verifyService.Verify(ctx, &request)

	c.JSON(http.StatusOK, response)
}
