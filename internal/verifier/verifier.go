package verifier

import (
	"context"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/pkg/errors"
)

// Verifier is the interface that all verifiers must implement
type Verifier interface {
	// Verify verifies a payment request
	Verify(ctx context.Context, request *models.VerifyRequest) VerificationResult
	// Type returns the verification step type
	Type() VerificationStep
	// Order returns the order in which this verifier should be executed
	Order() int
}

// VerificationResult represents the result of a verification process
type VerificationResult struct {
	IsValid           bool
	VerificationError *errors.X402Error
	ErrorMessage      string
}

// OK creates a successful verification result
func OK() VerificationResult {
	return VerificationResult{
		IsValid:           true,
		VerificationError: nil,
		ErrorMessage:      "",
	}
}

// Fail creates a failed verification result
func Fail(errorCode errors.X402Error, errorMessage string) VerificationResult {
	return VerificationResult{
		IsValid:           false,
		VerificationError: &errorCode,
		ErrorMessage:      errorMessage,
	}
}

// VerificationStep represents a step in the verification process
type VerificationStep string

const (
	// StepGlobalVerifier verifies that the request is globally valid
	StepGlobalVerifier VerificationStep = "GLOBAL_VERIFIER"
	// StepSchemeExists verifies the scheme
	StepSchemeExists VerificationStep = "SCHEME_EXISTS"
	// StepPaymentContextForExactScheme verifies payment context for exact scheme
	StepPaymentContextForExactScheme VerificationStep = "PAYMENT_CONTEXT_FOR_EXACT_SCHEME"
	// StepSignatureForExactScheme verifies signature for exact scheme
	StepSignatureForExactScheme VerificationStep = "SIGNATURE_FOR_EXACT_SCHEME"
	// StepPaymentAddressForExactScheme verifies payment address for exact scheme
	StepPaymentAddressForExactScheme VerificationStep = "PAYMENT_ADDRESS_FOR_EXACT_SCHEME"
	// StepDeadlinesForExactScheme checks deadlines for exact scheme
	StepDeadlinesForExactScheme VerificationStep = "DEADLINES_FOR_EXACT_SCHEME"
	// StepUserBalanceForExactScheme checks user balance for exact scheme
	StepUserBalanceForExactScheme VerificationStep = "USER_BALANCE_FOR_EXACT_SCHEME"
	// StepPaymentValueForExactScheme verifies payment value for exact scheme
	StepPaymentValueForExactScheme VerificationStep = "PAYMENT_VALUE_FOR_EXACT_SCHEME"
)

// String returns the string representation of the verification step
func (s VerificationStep) String() string {
	return string(s)
}
