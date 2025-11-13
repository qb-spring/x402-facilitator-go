package models

import "encoding/json"

// VerifyRequest represents a payment verification request
type VerifyRequest struct {
	X402Version         int                 `json:"x402Version" binding:"required"`
	PaymentPayload      PaymentPayload      `json:"paymentPayload" binding:"required"`
	PaymentRequirements PaymentRequirements `json:"paymentRequirements" binding:"required"`
}

// PaymentPayload represents the payment payload structure
type PaymentPayload struct {
	X402Version int     `json:"x402Version" binding:"required"`
	Scheme      string  `json:"scheme" binding:"required"`
	Network     string  `json:"network" binding:"required"`
	Payload     Payload `json:"payload" binding:"required"`
}

type Payload struct {
	Signature     string        `json:"signature" binding:"required,len=132,startswith=0x"`
	Authorization Authorization `json:"authorization" binding:"required"`
}

// Authorization represents the authorization data
type Authorization struct {
	From        string `json:"from" binding:"required,len=42,startswith=0x"`
	To          string `json:"to" binding:"required,len=42,startswith=0x"`
	Value       string `json:"value" binding:"required,numeric"`
	ValidAfter  string `json:"validAfter" binding:"required,numeric"`
	ValidBefore string `json:"validBefore" binding:"required,numeric"`
	Nonce       string `json:"nonce" binding:"required,len=66,startswith=0x"`
}

// PaymentRequirements represents payment requirements
type PaymentRequirements struct {
	Scheme            string          `json:"scheme" binding:"required"`
	Network           string          `json:"network" binding:"required"`
	MaxAmountRequired string          `json:"maxAmountRequired" binding:"required,numeric"`
	Resource          string          `json:"resource" binding:"required"`
	Description       string          `json:"description,omitempty"`
	MimeType          string          `json:"mimeType,omitempty"`
	OutputSchema      json.RawMessage `json:"outputSchema,omitempty"`
	PayTo             string          `json:"payTo" binding:"required,len=42,startswith=0x"`
	MaxTimeoutSeconds int             `json:"maxTimeoutSeconds" binding:"required"`
	Asset             string          `json:"asset" binding:"required,len=42,startswith=0x"`
	Extra             Extra           `json:"extra,omitempty"`
}

type Extra struct {
	Name    string `json:"name" binding:"omitempty"`
	Version string `json:"version" binding:"omitempty"`
}

// VerifyResponse represents a verification response
type VerifyResponse struct {
	IsValid       bool   `json:"isValid"`
	InvalidReason string `json:"invalidReason,omitempty"`
	Payer         string `json:"payer"`
}

// SettleRequest represents a payment settlement request
type SettleRequest struct {
	X402Version         int                 `json:"x402Version" binding:"required"`
	PaymentPayload      PaymentPayload      `json:"paymentPayload" binding:"required"`
	PaymentRequirements PaymentRequirements `json:"paymentRequirements" binding:"required"`
}

// SettleResponse represents a settlement response
type SettleResponse struct {
	Success     bool    `json:"success"`
	ErrorReason string  `json:"errorReason,omitempty"`
	Transaction *string `json:"transaction,omitempty"`
	Network     string  `json:"network"`
	Payer       string  `json:"payer"`
}

// SupportedKind represents a supported payment kind
type SupportedKind struct {
	X402Version int16  `json:"x402Version"`
	Scheme      string `json:"scheme"`
	Network     string `json:"network"`
}

// SupportedResponse represents the supported schemes and networks
type SupportedResponse struct {
	Kinds []SupportedKind `json:"kinds"`
}
