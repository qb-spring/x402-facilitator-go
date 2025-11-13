package exact

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/internal/web3"
	"x402-facilitator-go/pkg/errors"

	"go.uber.org/zap"
)

const (
	// ExactScheme is the exact payment scheme
	ExactScheme = "exact"
)

// PaymentContextVerifier verifies the payment context for exact scheme
type PaymentContextVerifier struct {
	web3Client *web3.Client
	logger     *zap.Logger
}

// NewPaymentContextVerifier creates a new PaymentContextVerifier
func NewPaymentContextVerifier(logger *zap.Logger, web3Client *web3.Client) *PaymentContextVerifier {
	return &PaymentContextVerifier{
		logger:     logger,
		web3Client: web3Client,
	}
}

// Verify verifies the payment context
func (p *PaymentContextVerifier) Verify(ctx context.Context, request *models.VerifyRequest) verifier.VerificationResult {
	paymentRequirements := request.PaymentRequirements
	paymentPayload := request.PaymentPayload
	if request.X402Version != 1 {
		return verifier.Fail(
			errors.ErrorInvalidX402Version,
			fmt.Sprintf("Unsupported X402 protocol version: %d", request.X402Version),
		)
	}

	// Check if scheme is supported (currently only "exact" is supported)
	if paymentRequirements.Scheme != ExactScheme {
		return verifier.Fail(
			errors.ErrorUnsupportedScheme,
			fmt.Sprintf("Unsupported scheme: %s", paymentRequirements.Scheme),
		)
	}

	// Schemes must match
	if paymentPayload.Scheme != paymentRequirements.Scheme {
		return verifier.Fail(
			errors.ErrorUnsupportedScheme,
			fmt.Sprintf("Scheme mismatch: payment payload scheme '%s' does not match payment requirements scheme '%s'",
				paymentPayload.Scheme, paymentRequirements.Scheme),
		)
	}

	// Check that the network in payment requirements is supported by trying to obtain its chain ID.
	_, err := p.web3Client.GetChainID(paymentRequirements.Network)
	if err != nil {
		return verifier.Fail(
			errors.ErrorInvalidNetwork,
			fmt.Sprintf("Network not supported: '%s'", paymentRequirements.Network),
		)
	}
	// Networks must match
	if paymentPayload.Network != paymentRequirements.Network {
		return verifier.Fail(
			errors.ErrorInvalidNetwork,
			fmt.Sprintf("Network mismatch: payment payload network '%s' does not match payment requirements network '%s'",
				paymentPayload.Network, paymentRequirements.Network),
		)
	}

	authorization := paymentPayload.Payload.Authorization

	if !strings.EqualFold(authorization.To, paymentRequirements.PayTo) {
		return verifier.Fail(
			errors.ErrorInvalidExactEVMPayloadRecipientMismatch,
			fmt.Sprintf(
				"Recipient mismatch: authorization.to '%s' does not match payment requirements payTo '%s'",
				authorization.To,
				paymentRequirements.PayTo,
			),
		)
	}

	// Validate numeric consistency between requirements and payload
	maxAmountRequired, _ := new(big.Int).SetString(paymentRequirements.MaxAmountRequired, 10)
	value, _ := new(big.Int).SetString(authorization.Value, 10)
	if value.Cmp(maxAmountRequired) < 0 {
		return verifier.Fail(
			errors.ErrorInvalidExactEVMPayloadAuthorizationValue,
			fmt.Sprintf(
				"Payment value is less than the required maximum amount (%s < %s)",
				authorization.Value,
				paymentRequirements.MaxAmountRequired,
			),
		)
	}

	// // Validate validAfter and validBefore (decimal strings, validBefore > validAfter)
	// validAfter, _ := new(big.Int).SetString(authorization.ValidAfter, 10)
	// validBefore, _ := new(big.Int).SetString(authorization.ValidBefore, 10)
	// // Validate current time falls within [validAfter, validBefore]
	// now := big.NewInt(time.Now().Unix())
	// if now.Cmp(validAfter) <= 0 {
	// 	return verifier.Fail(
	// 		errors.ErrorInvalidExactEVMPayloadAuthorizationValidAfter,
	// 		fmt.Sprintf("Authorization not yet valid: now=%d <= validAfter=%s", now, authorization.ValidAfter),
	// 	)
	// }
	// if now.Cmp(validBefore) >= 0 {
	// 	return verifier.Fail(
	// 		errors.ErrorInvalidExactEVMPayloadAuthorizationValidBefore,
	// 		fmt.Sprintf("Authorization expired: now=%d >= validBefore=%s", now, authorization.ValidBefore),
	// 	)
	// }

	return verifier.OK()
}

// Type returns the verification step type
func (p *PaymentContextVerifier) Type() verifier.VerificationStep {
	return verifier.StepPaymentContextForExactScheme
}

// Order returns the order in which this verifier should be executed
func (p *PaymentContextVerifier) Order() int {
	return 2
}
