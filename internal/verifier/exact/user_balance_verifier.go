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

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"go.uber.org/zap"
)

// UserBalanceVerifier verifies that the user has sufficient balance
type UserBalanceVerifier struct {
	logger     *zap.Logger
	web3Client *web3.Client
}

// NewUserBalanceVerifier creates a new UserBalanceVerifier
func NewUserBalanceVerifier(logger *zap.Logger, web3Client *web3.Client) *UserBalanceVerifier {
	return &UserBalanceVerifier{
		logger:     logger,
		web3Client: web3Client,
	}
}

// Verify verifies that the user has sufficient balance
func (u *UserBalanceVerifier) Verify(ctx context.Context, request *models.VerifyRequest) verifier.VerificationResult {
	ethCli, err := u.web3Client.GetClient(request.PaymentRequirements.Network)
	if err != nil {
		return verifier.Fail(
			errors.ErrorInvalidNetwork,
			fmt.Sprintf("Failed to get client for network: %v", err),
		)
	}

	contractAddr := common.HexToAddress(request.PaymentRequirements.Asset)
	userAddr := common.HexToAddress(request.PaymentPayload.Payload.Authorization.From)

	// Get user balance using ERC20 balanceOf method
	balance, err := u.getBalance(ctx, ethCli, contractAddr, userAddr)
	if err != nil {
		return verifier.Fail(
			errors.ErrorUnknown,
			fmt.Sprintf("Failed to get user balance: %v", err),
		)
	}

	// Parse the required value
	requiredValue, _ := new(big.Int).SetString(request.PaymentPayload.Payload.Authorization.Value, 10)

	// Check if balance is sufficient
	if balance.Cmp(requiredValue) < 0 {
		return verifier.Fail(
			errors.ErrorInsufficientFunds,
			fmt.Sprintf("Insufficient balance: user has %s, required %s", balance.String(), requiredValue.String()),
		)
	}

	return verifier.OK()
}

// getBalance retrieves the ERC20 token balance for an address
func (u *UserBalanceVerifier) getBalance(ctx context.Context, client bind.ContractCaller, contractAddr, userAddr common.Address) (*big.Int, error) {
	// ERC20 balanceOf function signature: balanceOf(address) returns (uint256)
	balanceOfABI := `[{"constant":true,"inputs":[{"name":"_owner","type":"address"}],"name":"balanceOf","outputs":[{"name":"balance","type":"uint256"}],"type":"function"}]`

	parsedABI, err := abi.JSON(strings.NewReader(balanceOfABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}

	// Create a bound contract
	boundContract := bind.NewBoundContract(contractAddr, parsedABI, client, nil, nil)

	// Call the contract
	var result []interface{}
	callOpts := &bind.CallOpts{Context: ctx}
	err = boundContract.Call(callOpts, &result, "balanceOf", userAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to call balanceOf: %w", err)
	}

	// Unpack the result
	balance, ok := result[0].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to unpack balance result")
	}
	return balance, nil
}

// Type returns the verification step type
func (u *UserBalanceVerifier) Type() verifier.VerificationStep {
	return verifier.StepUserBalanceForExactScheme
}

// Order returns the order in which this verifier should be executed
func (u *UserBalanceVerifier) Order() int {
	return 5
}
