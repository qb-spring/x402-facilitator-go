package service

import (
	"context"
	"math/big"
	"x402-facilitator-go/internal/models"
	"x402-facilitator-go/internal/web3"
	"x402-facilitator-go/internal/web3/contract"
	"x402-facilitator-go/pkg/errors"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"go.uber.org/zap"
)

// SettleService handles payment settlement
type SettleService struct {
	verifyService *VerifyService
	web3Client    *web3.Client
	privateKey    string
	logger        *zap.Logger
}

// NewSettleService creates a new SettleService
func NewSettleService(
	verifyService *VerifyService,
	web3Client *web3.Client,
	privateKey string,
	logger *zap.Logger,
) *SettleService {
	return &SettleService{
		verifyService: verifyService,
		web3Client:    web3Client,
		privateKey:    privateKey,
		logger:        logger,
	}
}

// Settle settles a payment request
func (s *SettleService) Settle(ctx context.Context, request *models.SettleRequest) *models.SettleResponse {
	// Verify the request first
	verifyRequest := &models.VerifyRequest{
		X402Version:         request.X402Version,
		PaymentPayload:      request.PaymentPayload,
		PaymentRequirements: request.PaymentRequirements,
	}

	verifyResponse := s.verifyService.Verify(ctx, verifyRequest)
	if !verifyResponse.IsValid {
		return &models.SettleResponse{
			Success:     false,
			Network:     request.PaymentRequirements.Network,
			ErrorReason: verifyResponse.InvalidReason,
			Payer:       verifyResponse.Payer,
		}
	}

	networkStr := request.PaymentRequirements.Network
	payer := verifyResponse.Payer
	auth := request.PaymentPayload.Payload.Authorization
	value, _ := new(big.Int).SetString(auth.Value, 10)
	validAfter, _ := new(big.Int).SetString(auth.ValidAfter, 10)
	validBefore, _ := new(big.Int).SetString(auth.ValidBefore, 10)
	nonceBytes := common.HexToHash(auth.Nonce)
	signatureBytes := common.FromHex(request.PaymentPayload.Payload.Signature)
	contractAddress := common.HexToAddress(request.PaymentRequirements.Asset)
	fromAddr := common.HexToAddress(auth.From)
	toAddr := common.HexToAddress(request.PaymentRequirements.PayTo)

	client, _ := s.web3Client.GetClient(networkStr)
	chainID, _ := s.web3Client.GetChainID(networkStr)

	// Parse private key and create transactor
	privateKey, err := crypto.HexToECDSA(s.privateKey)
	if err != nil {
		s.logger.Error("Invalid facilitator private key",
			zap.String("network", networkStr),
			zap.String("payer", payer),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorUnknown.Code(),
			Payer:       payer,
		}
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		s.logger.Error("Failed to create transactor with chain ID",
			zap.Error(err),
			zap.String("network", networkStr),
			zap.String("payer", payer),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorUnknown.Code(),
			Payer:       payer,
		}
	}
	transactOpts.Context = ctx

	// Create token contract instance
	tokenContract, err := contract.NewEIP3009Token(contractAddress, client)
	if err != nil {
		s.logger.Error("Failed to init EIP3009 token contract",
			zap.Error(err),
			zap.String("network", networkStr),
			zap.String("payer", payer),
			zap.String("contract", contractAddress.Hex()),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorUnknown.Code(),
			Payer:       payer,
		}
	}

	// Execute transfer
	tx, err := tokenContract.TransferWithAuthorization(
		transactOpts,
		fromAddr,
		toAddr,
		value,
		validAfter,
		validBefore,
		nonceBytes,
		signatureBytes,
	)
	if err != nil {
		s.logger.Warn("Token transferWithAuthorization reverted",
			zap.Error(err),
			zap.String("network", networkStr),
			zap.String("payer", payer),
			zap.String("contract", contractAddress.Hex()),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorInvalidTransactionState.Code(),
			Payer:       payer,
		}
	}

	// Wait for confirmation
	s.logger.Info("Transaction sent, waiting for confirmation",
		zap.String("txHash", tx.Hash().Hex()),
		zap.String("network", networkStr),
		zap.String("payer", payer),
	)

	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		s.logger.Warn("Failed while waiting for tx receipt",
			zap.String("txHash", tx.Hash().Hex()),
			zap.Error(err),
			zap.String("network", networkStr),
			zap.String("payer", payer),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorUnexpectedSettle.Code(),
			Payer:       payer,
		}
	}

	if receipt.Status == types.ReceiptStatusFailed {
		s.logger.Warn("Settlement transaction failed on-chain",
			zap.String("txHash", tx.Hash().Hex()),
			zap.Uint64("blockNumber", receipt.BlockNumber.Uint64()),
			zap.String("network", networkStr),
			zap.String("payer", payer),
		)
		return &models.SettleResponse{
			Success:     false,
			Network:     networkStr,
			ErrorReason: errors.ErrorInvalidTransactionState.Code(),
			Payer:       payer,
		}
	}

	txHash := tx.Hash().Hex()
	s.logger.Info("Settlement transaction confirmed",
		zap.String("txHash", txHash),
		zap.Uint64("blockNumber", receipt.BlockNumber.Uint64()),
		zap.String("network", networkStr),
		zap.String("payer", payer),
	)

	return &models.SettleResponse{
		Success:     true,
		Network:     networkStr,
		Transaction: &txHash,
		Payer:       payer,
	}
}
