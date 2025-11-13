package exact

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/pkg/errors"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// GlobalVerifier verifies that the request is globally valid
type GlobalVerifier struct {
	validator *validator.Validate
	logger    *zap.Logger
}

// NewGlobalVerifier creates a new GlobalVerifier
func NewGlobalVerifier(logger *zap.Logger) *GlobalVerifier {
	v := validator.New()

	// Register custom tag name function to use json tag names
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return &GlobalVerifier{
		validator: v,
		logger:    logger,
	}
}

// Verify verifies the request
func (g *GlobalVerifier) Verify(ctx context.Context, request *models.VerifyRequest) verifier.VerificationResult {
	// Validate the request
	if err := g.validator.Struct(request); err != nil {
		// Get the first validation error
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range validationErrors {
				errorMsg := g.getErrorMessage(fieldError)
				return verifier.Fail(errors.ErrorInvalidPayload, errorMsg)
			}
		} else {
			return verifier.Fail(errors.ErrorInvalidPayload, err.Error())
		}
	}

	return verifier.OK()
}

// Type returns the verification step type
func (g *GlobalVerifier) Type() verifier.VerificationStep {
	return verifier.StepGlobalVerifier
}

// Order returns the order in which this verifier should be executed
func (g *GlobalVerifier) Order() int {
	return 1
}

// getErrorMessage generates a human-readable error message from a validation error
func (g *GlobalVerifier) getErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()

	switch tag {
	case "required":
		return fmt.Sprintf("Field '%s' is required", field)
	case "min":
		return fmt.Sprintf("Field '%s' must be at least %s", field, param)
	case "max":
		return fmt.Sprintf("Field '%s' must be at most %s", field, param)
	case "email":
		return fmt.Sprintf("Field '%s' must be a valid email address", field)
	case "url":
		return fmt.Sprintf("Field '%s' must be a valid URL", field)
	default:
		return fmt.Sprintf("Field '%s' failed validation: %s", field, tag)
	}
}
