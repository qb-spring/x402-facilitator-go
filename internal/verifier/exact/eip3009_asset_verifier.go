package exact

import (
	"context"
	"fmt"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/internal/web3"
	"x402-facilitator-go/internal/web3/contract"
	"x402-facilitator-go/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// EIP3009AssetVerifier checks the asset supports EIP-3009 transferWithAuthorization
type EIP3009AssetVerifier struct {
	logger     *zap.Logger
	web3Client *web3.Client
}

// NewEIP3009AssetVerifier creates a new EIP3009AssetVerifier
func NewEIP3009AssetVerifier(logger *zap.Logger, web3Client *web3.Client) *EIP3009AssetVerifier {
	return &EIP3009AssetVerifier{
		logger:     logger,
		web3Client: web3Client,
	}
}

// Verify verifies the asset contract supports EIP-3009 by probing authorizationState
func (v *EIP3009AssetVerifier) Verify(ctx context.Context, request *models.VerifyRequest) verifier.VerificationResult {
	ethCli, _ := v.web3Client.GetClient(request.PaymentRequirements.Network)

	contractAddr := common.HexToAddress(request.PaymentRequirements.Asset)

	code, err := ethCli.CodeAt(ctx, contractAddr, nil)
	if err != nil {
		return verifier.Fail(
			errors.ErrorUnknown,
			fmt.Sprintf("Failed to fetch asset bytecode: %v", err),
		)
	}
	if len(code) == 0 {
		return verifier.Fail(
			errors.ErrorInvalidPayload,
			"Asset address is not a contract",
		)
	}

	tokenContract, err := contract.NewEIP3009Token(contractAddr, ethCli)
	if err != nil {
		return verifier.Fail(errors.ErrorUnknown, fmt.Sprintf("Failed to create EIP3009Token contract instance: %v", err))
	}

	// Prepare a harmless staticcall to check method existence
	zeroAddr := common.Address{} // 0x000...00
	var zeroNonce [32]byte

	_, err = tokenContract.AuthorizationState(&bind.CallOpts{Context: ctx}, zeroAddr, zeroNonce)
	if err != nil {
		// Method missing or reverted → not EIP-3009
		return verifier.Fail(
			errors.ErrorInvalidPayload,
			"Asset does not support EIP-3009 transferWithAuthorization (authorizationState missing)",
		)
	}

	return verifier.OK()
}

// Type returns the verification step type
func (v *EIP3009AssetVerifier) Type() verifier.VerificationStep {
	return verifier.StepPaymentAddressForExactScheme
}

// Order returns the order in which this verifier should be executed
func (v *EIP3009AssetVerifier) Order() int {
	return 3
}
