package eip712

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// DomainSeparatorParams represents the parameters for EIP-712 domain separator
type DomainSeparatorParams struct {
	Name              string
	Version           string
	ChainID           *big.Int
	VerifyingContract common.Address
}

// ComputeDomainSeparator computes the EIP-712 domain separator hash
func ComputeDomainSeparator(params DomainSeparatorParams) []byte {
	domainTypehash := crypto.Keccak256Hash([]byte("EIP712Domain(string name,string version,uint256 chainId,address verifyingContract)"))

	domainSeparatorHash := crypto.Keccak256(
		domainTypehash.Bytes(),
		crypto.Keccak256([]byte(params.Name)),
		crypto.Keccak256([]byte(params.Version)),
		common.LeftPadBytes(params.ChainID.Bytes(), 32),
		common.LeftPadBytes(params.VerifyingContract.Bytes(), 32),
	)

	return domainSeparatorHash
}

// ComputeMessageHash computes the EIP-712 message hash given the typehash and encoded data
func ComputeMessageHash(typehash []byte, encodedData []byte) []byte {
	return crypto.Keccak256(typehash, encodedData)
}

// ComputeEIP712Hash computes the final EIP-712 hash from domain separator and message hash
func ComputeEIP712Hash(domainSeparatorHash []byte, messageHash []byte) []byte {
	eip712Prefix := []byte("\x19\x01")
	return crypto.Keccak256(eip712Prefix, domainSeparatorHash, messageHash)
}

// VerifySignature verifies an EIP-712 signature and recovers the signer address
func VerifySignature(hashBytes []byte, signatureHex string, expectedAddress common.Address) (bool, common.Address, error) {
	sigBytes := common.FromHex(signatureHex)

	sigForRecovery := make([]byte, len(sigBytes))
	copy(sigForRecovery, sigBytes)
	// Ethereum signatures use v = 27/28 or v = 0/1, normalize if necessary
	if sigForRecovery[64] >= 27 {
		sigForRecovery[64] -= 27
	}

	pubKey, err := crypto.SigToPub(hashBytes, sigForRecovery)
	if err != nil {
		return false, common.Address{}, fmt.Errorf("failed to recover public key: %w", err)
	}

	signerAddress := crypto.PubkeyToAddress(*pubKey)

	// Compare addresses
	if signerAddress != expectedAddress {
		return false, signerAddress, fmt.Errorf(
			"signature mismatch: expected %s, got %s",
			expectedAddress.Hex(),
			signerAddress.Hex(),
		)
	}

	return true, signerAddress, nil
}
