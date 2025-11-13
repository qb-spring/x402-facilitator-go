package exact

import (
	"context"
	"fmt"
	"math/big"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/util/eip3009"
	"x402-facilitator-go/internal/util/eip712"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/internal/web3"
	"x402-facilitator-go/pkg/errors"

	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// var transferWithAuthorizationTypehash = crypto.Keccak256Hash([]byte("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))

// SignatureVerifier verifies the signature for exact scheme
type SignatureVerifier struct {
	logger     *zap.Logger
	web3Client *web3.Client
}

// NewSignatureVerifier creates a new SignatureVerifier
func NewSignatureVerifier(logger *zap.Logger, web3Client *web3.Client) *SignatureVerifier {
	return &SignatureVerifier{
		logger:     logger,
		web3Client: web3Client,
	}
}

// Verify verifies the signature
func (s *SignatureVerifier) Verify(ctx context.Context, request *models.VerifyRequest) verifier.VerificationResult {
	// Resolve chain ID first to validate network
	chainId, _ := s.web3Client.GetChainID(request.PaymentPayload.Network)
	// Compute EIP-712 hash
	hashBytes := s.computeTransferWithAuthorizationHash(request, chainId)

	// Get the exact scheme payload
	exactPayload := &request.PaymentPayload.Payload
	// Verify signature using eip712 utility
	expectedAddress := common.HexToAddress(exactPayload.Authorization.From)
	isValid, signerAddress, err := eip712.VerifySignature(hashBytes, exactPayload.Signature, expectedAddress)
	if err != nil {
		return verifier.Fail(
			errors.ErrorInvalidExactEVMPayloadSignature,
			fmt.Sprintf("Signature verification failed: %v", err),
		)
	}

	if !isValid {
		return verifier.Fail(
			errors.ErrorInvalidExactEVMPayloadSignature,
			fmt.Sprintf("Signature mismatch: expected %s, got %s",
				exactPayload.Authorization.From,
				signerAddress.Hex()),
		)
	}
	return verifier.OK()
}

func (s *SignatureVerifier) computeTransferWithAuthorizationHash(req *models.VerifyRequest, chainId *big.Int) []byte {
	params := eip3009.TransferWithAuthorizationParams{
		ChainId:           chainId,
		VerifyingContract: req.PaymentRequirements.Asset,
		DomainName:        req.PaymentRequirements.Extra.Name,
		DomainVersion:     req.PaymentRequirements.Extra.Version,
		From:              req.PaymentPayload.Payload.Authorization.From,
		To:                req.PaymentPayload.Payload.Authorization.To,
		Value:             req.PaymentPayload.Payload.Authorization.Value,
		ValidAfter:        req.PaymentPayload.Payload.Authorization.ValidAfter,
		ValidBefore:       req.PaymentPayload.Payload.Authorization.ValidBefore,
		Nonce:             req.PaymentPayload.Payload.Authorization.Nonce,
	}

	return eip3009.ComputeTransferWithAuthorizationHash(params)
}

// Type returns the verification step type
func (s *SignatureVerifier) Type() verifier.VerificationStep {
	return verifier.StepSignatureForExactScheme
}

// Order returns the order in which this verifier should be executed
func (s *SignatureVerifier) Order() int {
	return 4
}
