package service

import (
	"context"
	"sort"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/pkg/errors"

	"go.uber.org/zap"
)

// VerifyService handles payment verification
type VerifyService struct {
	verifiers []verifier.Verifier
	logger    *zap.Logger
}

// NewVerifyService creates a new VerifyService
func NewVerifyService(verifiers []verifier.Verifier, logger *zap.Logger) *VerifyService {
	// Sort verifiers by order
	sort.Slice(verifiers, func(i, j int) bool {
		return verifiers[i].Order() < verifiers[j].Order()
	})

	logger.Debug("Verify service initialized",
		zap.Int("verifierCount", len(verifiers)),
	)

	return &VerifyService{
		verifiers: verifiers,
		logger:    logger,
	}
}

// Verify verifies a payment request
func (s *VerifyService) Verify(ctx context.Context, request *models.VerifyRequest) *models.VerifyResponse {
	// Run all verifiers in order, return the first failure if any
	payer := request.PaymentPayload.Payload.Authorization.From
	for _, v := range s.verifiers {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			s.logger.Info("Verification cancelled",
				zap.String("verifier", v.Type().String()),
				zap.Error(ctx.Err()),
				zap.String("network", request.PaymentRequirements.Network),
				zap.String("payer", payer),
			)
			return &models.VerifyResponse{
				IsValid:       false,
				InvalidReason: errors.ErrorUnknown.Code(),
				Payer:         payer,
			}
		default:
		}

		s.logger.Debug("Running verification",
			zap.String("verifier", v.Type().String()),
			zap.String("network", request.PaymentRequirements.Network),
			zap.String("payer", payer),
		)

		result := v.Verify(ctx, request)
		if !result.IsValid {
			s.logger.Warn("Verification failed",
				zap.String("verifier", v.Type().String()),
				zap.String("error", result.ErrorMessage),
				zap.String("network", request.PaymentRequirements.Network),
				zap.String("payer", payer),
			)
			return &models.VerifyResponse{
				IsValid:       false,
				InvalidReason: result.VerificationError.Code(),
				Payer:         payer,
			}
		}

		s.logger.Debug("Verification passed",
			zap.String("verifier", v.Type().String()),
			zap.String("network", request.PaymentRequirements.Network),
			zap.String("payer", payer),
		)
	}

	// All verifiers passed
	s.logger.Debug("All verifiers passed",
		zap.String("network", request.PaymentRequirements.Network),
		zap.String("payer", payer),
	)
	return &models.VerifyResponse{
		IsValid: true,
		Payer:   payer,
	}
}
