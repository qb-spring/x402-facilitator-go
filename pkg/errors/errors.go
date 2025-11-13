package errors

// X402Error represents an X402 error code
type X402Error string

const (

	// invalid_scheme

	// Blockchain transaction failed or was rejected
	ErrorInvalidTransactionState X402Error = "INVALID_TRANSACTION_STATE"
	// Unexpected error occurred during payment verification
	ErrorUnexpectedVerify X402Error = "UNEXPECTED_VERIFY_ERROR"
	// Unexpected error occurred during payment settlement
	ErrorUnexpectedSettle X402Error = "UNEXPECTED_SETTLE_ERROR"

	// Verify errors
	// ErrorUnknown represents an unknown error
	ErrorUnknown X402Error = "UNKNOWN"
	// Protocol version is not supported
	ErrorInvalidX402Version X402Error = "INVALID_X402_VERSION"
	// ErrorInvalidPayload represents an invalid payload error
	ErrorInvalidPayload X402Error = "INVALID_PAYLOAD"
	// ErrorUnsupportedScheme represents an unsupported scheme error
	ErrorUnsupportedScheme X402Error = "UNSUPPORTED_SCHEME"
	// ErrorInvalidNetwork represents an invalid network error
	ErrorInvalidNetwork X402Error = "INVALID_NETWORK"
	// ErrorInvalidExactEVMPayloadSignature represents an invalid signature error
	ErrorInvalidExactEVMPayloadSignature X402Error = "INVALID_EXACT_EVM_PAYLOAD_SIGNATURE"
	// ErrorInvalidExactEVMPayloadAuthorizationValue represents an invalid authorization value error
	ErrorInvalidExactEVMPayloadAuthorizationValue X402Error = "INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALUE"
	// ErrorInvalidExactEVMPayloadRecipientMismatch represents a recipient mismatch error
	ErrorInvalidExactEVMPayloadRecipientMismatch X402Error = "INVALID_EXACT_EVM_PAYLOAD_RECIPIENT_MISMATCH"
	// ErrorInvalidExactEVMPayloadAuthorizationValidAfter represents an invalid valid after error
	ErrorInvalidExactEVMPayloadAuthorizationValidAfter X402Error = "INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_AFTER"
	// ErrorInvalidExactEVMPayloadAuthorizationValidBefore represents an invalid valid before error
	ErrorInvalidExactEVMPayloadAuthorizationValidBefore X402Error = "INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_BEFORE"
	// ErrorInsufficientFunds represents an insufficient funds error
	ErrorInsufficientFunds X402Error = "INSUFFICIENT_FUNDS"
)

// Code returns the error code string
func (e X402Error) Code() string {
	return string(e)
}
