// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// EIP3009TokenMetaData contains all meta data concerning the EIP3009Token contract.
var EIP3009TokenMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validAfter\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validBefore\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"transferWithAuthorization\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"authorizer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"nonce\",\"type\":\"bytes32\"}],\"name\":\"authorizationState\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EIP3009TokenABI is the input ABI used to generate the binding from.
// Deprecated: Use EIP3009TokenMetaData.ABI instead.
var EIP3009TokenABI = EIP3009TokenMetaData.ABI

// EIP3009Token is an auto generated Go binding around an Ethereum contract.
type EIP3009Token struct {
	EIP3009TokenCaller     // Read-only binding to the contract
	EIP3009TokenTransactor // Write-only binding to the contract
	EIP3009TokenFilterer   // Log filterer for contract events
}

// EIP3009TokenCaller is an auto generated read-only Go binding around an Ethereum contract.
type EIP3009TokenCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP3009TokenTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EIP3009TokenTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP3009TokenFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EIP3009TokenFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EIP3009TokenSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EIP3009TokenSession struct {
	Contract     *EIP3009Token     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EIP3009TokenCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EIP3009TokenCallerSession struct {
	Contract *EIP3009TokenCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// EIP3009TokenTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EIP3009TokenTransactorSession struct {
	Contract     *EIP3009TokenTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// EIP3009TokenRaw is an auto generated low-level Go binding around an Ethereum contract.
type EIP3009TokenRaw struct {
	Contract *EIP3009Token // Generic contract binding to access the raw methods on
}

// EIP3009TokenCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EIP3009TokenCallerRaw struct {
	Contract *EIP3009TokenCaller // Generic read-only contract binding to access the raw methods on
}

// EIP3009TokenTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EIP3009TokenTransactorRaw struct {
	Contract *EIP3009TokenTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEIP3009Token creates a new instance of EIP3009Token, bound to a specific deployed contract.
func NewEIP3009Token(address common.Address, backend bind.ContractBackend) (*EIP3009Token, error) {
	contract, err := bindEIP3009Token(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EIP3009Token{EIP3009TokenCaller: EIP3009TokenCaller{contract: contract}, EIP3009TokenTransactor: EIP3009TokenTransactor{contract: contract}, EIP3009TokenFilterer: EIP3009TokenFilterer{contract: contract}}, nil
}

// NewEIP3009TokenCaller creates a new read-only instance of EIP3009Token, bound to a specific deployed contract.
func NewEIP3009TokenCaller(address common.Address, caller bind.ContractCaller) (*EIP3009TokenCaller, error) {
	contract, err := bindEIP3009Token(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EIP3009TokenCaller{contract: contract}, nil
}

// NewEIP3009TokenTransactor creates a new write-only instance of EIP3009Token, bound to a specific deployed contract.
func NewEIP3009TokenTransactor(address common.Address, transactor bind.ContractTransactor) (*EIP3009TokenTransactor, error) {
	contract, err := bindEIP3009Token(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EIP3009TokenTransactor{contract: contract}, nil
}

// NewEIP3009TokenFilterer creates a new log filterer instance of EIP3009Token, bound to a specific deployed contract.
func NewEIP3009TokenFilterer(address common.Address, filterer bind.ContractFilterer) (*EIP3009TokenFilterer, error) {
	contract, err := bindEIP3009Token(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EIP3009TokenFilterer{contract: contract}, nil
}

// bindEIP3009Token binds a generic wrapper to an already deployed contract.
func bindEIP3009Token(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := EIP3009TokenMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP3009Token *EIP3009TokenRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP3009Token.Contract.EIP3009TokenCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP3009Token *EIP3009TokenRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP3009Token.Contract.EIP3009TokenTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP3009Token *EIP3009TokenRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP3009Token.Contract.EIP3009TokenTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EIP3009Token *EIP3009TokenCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EIP3009Token.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EIP3009Token *EIP3009TokenTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EIP3009Token.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EIP3009Token *EIP3009TokenTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EIP3009Token.Contract.contract.Transact(opts, method, params...)
}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_EIP3009Token *EIP3009TokenCaller) AuthorizationState(opts *bind.CallOpts, authorizer common.Address, nonce [32]byte) (bool, error) {
	var out []interface{}
	err := _EIP3009Token.contract.Call(opts, &out, "authorizationState", authorizer, nonce)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_EIP3009Token *EIP3009TokenSession) AuthorizationState(authorizer common.Address, nonce [32]byte) (bool, error) {
	return _EIP3009Token.Contract.AuthorizationState(&_EIP3009Token.CallOpts, authorizer, nonce)
}

// AuthorizationState is a free data retrieval call binding the contract method 0xe94a0102.
//
// Solidity: function authorizationState(address authorizer, bytes32 nonce) view returns(bool)
func (_EIP3009Token *EIP3009TokenCallerSession) AuthorizationState(authorizer common.Address, nonce [32]byte) (bool, error) {
	return _EIP3009Token.Contract.AuthorizationState(&_EIP3009Token.CallOpts, authorizer, nonce)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_EIP3009Token *EIP3009TokenTransactor) TransferWithAuthorization(opts *bind.TransactOpts, from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _EIP3009Token.contract.Transact(opts, "transferWithAuthorization", from, to, value, validAfter, validBefore, nonce, signature)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_EIP3009Token *EIP3009TokenSession) TransferWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _EIP3009Token.Contract.TransferWithAuthorization(&_EIP3009Token.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}

// TransferWithAuthorization is a paid mutator transaction binding the contract method 0xcf092995.
//
// Solidity: function transferWithAuthorization(address from, address to, uint256 value, uint256 validAfter, uint256 validBefore, bytes32 nonce, bytes signature) returns()
func (_EIP3009Token *EIP3009TokenTransactorSession) TransferWithAuthorization(from common.Address, to common.Address, value *big.Int, validAfter *big.Int, validBefore *big.Int, nonce [32]byte, signature []byte) (*types.Transaction, error) {
	return _EIP3009Token.Contract.TransferWithAuthorization(&_EIP3009Token.TransactOpts, from, to, value, validAfter, validBefore, nonce, signature)
}
