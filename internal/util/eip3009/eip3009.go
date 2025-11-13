package eip3009

import (
	"math/big"
	"x402-facilitator-go/internal/util/eip712"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// Precomputed EIP-712 typehash for TransferWithAuthorization message
var transferWithAuthorizationTypehash = crypto.Keccak256Hash([]byte("TransferWithAuthorization(address from,address to,uint256 value,uint256 validAfter,uint256 validBefore,bytes32 nonce)"))

// TransferWithAuthorizationParams represents the parameters for computing EIP-3009 hash
type TransferWithAuthorizationParams struct {
	ChainId *big.Int
	// VerifyingContract is the contract address that will verify the signature
	VerifyingContract string
	// DomainName is the EIP-712 domain name (default: "USD Coin" if empty)
	DomainName string
	// DomainVersion is the EIP-712 domain version (default: "2" if empty)
	DomainVersion string
	// Authorization data
	From        string
	To          string
	Value       string
	ValidAfter  string
	ValidBefore string
	Nonce       string
}

// ComputeTransferWithAuthorizationHash computes the EIP-712 hash for the TransferWithAuthorization message
// This implements EIP-3009's transferWithAuthorization signature verification
func ComputeTransferWithAuthorizationHash(params TransferWithAuthorizationParams) []byte {
	// Compute domain separator using eip712 utility
	domainParams := eip712.DomainSeparatorParams{
		Name:              params.DomainName,
		Version:           params.DomainVersion,
		ChainID:           params.ChainId,
		VerifyingContract: common.HexToAddress(params.VerifyingContract),
	}
	domainSeparatorHash := eip712.ComputeDomainSeparator(domainParams)

	// Parse authorization values
	value, _ := new(big.Int).SetString(params.Value, 10)
	validAfter, _ := new(big.Int).SetString(params.ValidAfter, 10)
	validBefore, _ := new(big.Int).SetString(params.ValidBefore, 10)
	nonceHash := common.HexToHash(params.Nonce)
	fromAddr := common.HexToAddress(params.From)
	toAddr := common.HexToAddress(params.To)
	messageHash := crypto.Keccak256(
		transferWithAuthorizationTypehash.Bytes(),
		common.LeftPadBytes(fromAddr.Bytes(), 32),
		common.LeftPadBytes(toAddr.Bytes(), 32),
		common.LeftPadBytes(value.Bytes(), 32),
		common.LeftPadBytes(validAfter.Bytes(), 32),
		common.LeftPadBytes(validBefore.Bytes(), 32),
		nonceHash.Bytes(),
	)

	// Compute final EIP-712 hash using eip712 utility
	finalHash := eip712.ComputeEIP712Hash(domainSeparatorHash, messageHash)

	return finalHash
}
