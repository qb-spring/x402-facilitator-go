# X402 Facilitator

A Go-based implementation of the X402 payment protocol Facilitator service that supports verification and settlement of authorized transfer payments based on the EIP-3009 standard.

## Table of Contents

- [Project Structure](#project-structure)
- [About X402](#about-x402)
- [Features](#features)
- [Architecture](#architecture)
- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Development Guide](#development-guide)
- [Error Codes](#error-codes)
- [Security Considerations](#security-considerations)
- [Useful links](#useful-links)

## Project Structure

```
x402-facilitator/
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point, initializes services and routes
│
├── internal/
│   ├── config/
│   │   └── config.go                  # Configuration management, loads YAML config and environment variables
│   │
│   ├── handlers/
│   │   ├── verify_handler.go          # Verification request handler (POST /verify)
│   │   ├── settle_handler.go          # Settlement request handler (POST /settle)
│   │   └── supported_handler.go       # Supported networks/schemes query handler (GET /supported)
│   │
│   ├── middleware/
│   │   ├── cors.go                    # CORS cross-origin middleware
│   │   ├── logger.go                  # Request logging middleware
│   │   └── recovery.go                # Error recovery middleware
│   │
│   ├── models/
│   │   └── models.go                  # Data model definitions (request/response structs)
│   │
│   ├── service/
│   │   ├── verify_service.go          # Verification service, coordinates multiple verifiers
│   │   ├── settle_service.go          # Settlement service, executes on-chain token transfers
│   │   └── supported_service.go       # Supported networks/schemes query service
│   │
│   ├── util/
│   │   ├── eip3009/
│   │   │   └── eip3009.go             # EIP-3009 utility functions, calculates authorization hash
│   │   └── eip712/
│   │       └── eip712.go              # EIP-712 utility functions, signature verification
│   │
│   ├── verifier/
│   │   ├── verifier.go                # Verifier interface definition
│   │   └── exact/
│   │       ├── global_verifier.go              # Global verifier (Order: 1)
│   │       ├── payment_context_verifier.go     # Payment context verifier (Order: 2)
│   │       ├── eip3009_asset_verifier.go       # EIP-3009 asset verifier (Order: 3)
│   │       ├── signature_verifier.go           # Signature verifier (Order: 4)
│   │       └── user_balance_verifier.go       # User balance verifier (Order: 5)
│   │
│   └── web3/
│       ├── client.go                  # Web3 client management, supports multiple networks
│       └── contract/
│           └── EIP3009Token.go        # EIP-3009 contract ABI bindings
│
├── pkg/
│   └── errors/
│       └── errors.go                  # X402 error code definitions
│
├── config.yaml                         # Configuration file example
├── go.mod                              # Go module definition
├── go.sum                              # Dependency checksums
└── README.md                           # Project documentation
```

### Directory Descriptions

#### `cmd/server/`
Application entry point, responsible for:
- Loading configuration
- Initializing logging system
- Creating Web3 clients
- Registering verifiers
- Initializing services and handlers
- Setting up HTTP routes
- Graceful shutdown

#### `internal/config/`
Configuration management module:
- Loads configuration from YAML files
- Loads sensitive information (e.g., private keys) from environment variables
- Configuration validation

#### `internal/handlers/`
HTTP request handlers:
- `VerifyHandler`: Handles payment verification requests
- `SettleHandler`: Handles payment settlement requests
- `SupportedHandler`: Returns list of supported networks and schemes

#### `internal/middleware/`
HTTP middleware:
- `CORS`: Handles cross-origin requests
- `Logger`: Logs request information
- `Recovery`: Catches panics and returns error responses

#### `internal/models/`
Data model definitions:
- `VerifyRequest`: Verification request structure
- `SettleRequest`: Settlement request structure
- `PaymentPayload`: Payment payload
- `PaymentRequirements`: Payment requirements
- `Authorization`: Authorization information
- Response structures

#### `internal/service/`
Business logic layer:
- `VerifyService`: Coordinates multiple verifiers to execute verification in order
- `SettleService`: Executes on-chain token transfers
- `SupportedService`: Returns supported network configurations

#### `internal/verifier/`
Verifier module, implements chain verification:
- `Verifier` interface: Defines standard verifier interface
- `exact/`: Implements verifiers for "exact" payment scheme
  - Executes in order defined by `Order()` method
  - Any verifier failure immediately returns

#### `internal/util/`
Utility functions:
- `eip3009/`: EIP-3009 standard related utilities
- `eip712/`: EIP-712 structured data signature utilities

#### `internal/web3/`
Blockchain interaction layer:
- `Client`: Manages Ethereum clients for multiple networks
- `contract/`: Smart contract ABI bindings

#### `pkg/errors/`
Error code definitions:
- Defines all X402 protocol error codes
- Provides error code string conversion methods

## About X402

X402 is a decentralized payment protocol designed to provide standardized payment solutions for Web3 applications. The protocol allows users to complete payments through signature authorization without manually confirming transactions for each payment.

### X402 Core Concepts

1. **Payment Authorization**: Users create payment authorizations through EIP-712 signatures, including:
   - Payer address (From)
   - Payee address (To)
   - Payment amount (Value)
   - Validity period (ValidAfter, ValidBefore)
   - Nonce

2. **Payment Scheme**: Currently supports the `exact` scheme, which requires the payment amount to exactly match the authorization amount.

3. **Facilitator**: Responsible for verifying the validity of payment authorizations and executing on-chain settlements after verification passes.

4. **EIP-3009 Standard**: An extension standard based on ERC-20 that supports token transfers through authorization signatures without requiring users to manually confirm each time.

## Features

### 1. Multi-Network Support

Supports multiple EVM-compatible networks, currently configured with:
- Base Sepolia (testnet)
- Base Mainnet (mainnet)

Additional networks can be easily added through configuration files.

### 2. Multi-Layer Verification

Implements a complete verification chain executed in order:

1. **Global Verifier**: Validates request format and required fields
2. **Payment Context Verifier**: Validates protocol version, scheme, and network matching
3. **EIP-3009 Asset Verifier**: Validates whether token contracts support EIP-3009
4. **Signature Verifier**: Validates payment authorization signatures using EIP-712
5. **User Balance Verifier**: Validates whether user account balance is sufficient

### 3. Security Features

- EIP-712 structured data signature verification
- Private keys managed through environment variables, not stored in configuration files
- Complete error handling and logging
- CORS support

### 4. High Availability

- Graceful shutdown
- Context cancellation support
- Structured logging (JSON/Console format)
- Health check endpoint

## Architecture

### Layered Architecture

```
┌─────────────────────────────────────┐
│         HTTP Handlers              │
│  (verify, settle, supported)       │
└─────────────────────────────────────┘
              │
┌─────────────────────────────────────┐
│         Service Layer               │
│  (VerifyService, SettleService)     │
└─────────────────────────────────────┘
              │
┌─────────────────────────────────────┐
│         Verifier Layer              │
│  (Multiple verifiers executed in order) │
└─────────────────────────────────────┘
              │
┌─────────────────────────────────────┐
│         Web3 Client Layer           │
│  (Multi-network RPC client management) │
└─────────────────────────────────────┘
```

### Verifier Execution Order

Verifiers execute in the order returned by the `Order()` method:

1. **Order 1**: `GlobalVerifier` - Global format validation
2. **Order 2**: `PaymentContextVerifier` - Payment context validation
3. **Order 3**: `EIP3009AssetVerifier` - Asset contract validation
4. **Order 4**: `SignatureVerifier` - Signature validation
5. **Order 5**: `UserBalanceVerifier` - Balance validation

Any verifier failure immediately returns without continuing to subsequent verifiers.

### Data Flow

```
Client Request
    ↓
HTTP Handler (Parse request)
    ↓
Service Layer (Business logic)
    ↓
Verifier Chain (Chain verification)
    ↓
Web3 Client (Blockchain interaction)
    ↓
Return Response
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- Access to RPC nodes for EVM-compatible networks
- Facilitator private key (for executing settlement transactions)

### Installation

```bash
# Clone repository
git clone <repository-url>
cd x402-facilitator

# Install dependencies
go mod download
```

### Configuration

1. Copy and edit the configuration file:

```bash
cp config.yaml config.yaml.local
```

2. Edit `config.yaml.local` to configure server and network information.

3. Set environment variables:

```bash
export X402_FACILITATOR_PRIVATE_KEY="your_private_key_here"
```

### Running

```bash
# Use default configuration
go run cmd/server/main.go

# Or specify configuration file path
CONFIG_PATH=./config.yaml.local go run cmd/server/main.go
```

### Building

```bash
# Build binary
go build -o bin/x402-facilitator cmd/server/main.go

# Run
./bin/x402-facilitator
```

## Configuration

### Configuration File Structure

```yaml
server:
  host: "0.0.0.0"      # Server listen address
  port: 8081           # Server port

logging:
  level: "info"        # Log level: debug, info, warn, error
  format: "json"       # Log format: json, console

networks:
  networkInfos:
    - name: "base-sepolia"           # Network name (for API requests)
      rpcURL: "https://sepolia.base.org"  # RPC node URL
      chainId: 84532                 # Chain ID
      X402Version: 1                 # Supported X402 protocol version
      scheme: "exact"                # Supported payment scheme
```

### Environment Variables

- `X402_FACILITATOR_PRIVATE_KEY`: Facilitator private key (required)
- `CONFIG_PATH`: Configuration file path (optional)

### Configuration File Search Order

1. Path specified by `CONFIG_PATH` environment variable
2. `config.yaml` in current working directory
3. `config.yaml` in project root directory

## Development Guide

### Adding a New Verifier

1. Create a new verifier file in the `internal/verifier/exact/` directory
2. Implement the `verifier.Verifier` interface:
   ```go
   type Verifier interface {
       Verify(ctx context.Context, request *models.VerifyRequest) VerificationResult
       Type() VerificationStep
       Order() int
   }
   ```
3. Register the verifier in `cmd/server/main.go`:
   ```go
   verifiers := []verifier.Verifier{
       // ... existing verifiers
       exact.NewYourVerifier(logger, web3Client),
   }
   ```

### Adding a New Network

Add network configuration in `config.yaml`:

```yaml
networks:
  networkInfos:
    - name: "your-network"
      rpcURL: "https://your-rpc-url"
      chainId: 12345
      X402Version: 1
      scheme: "exact"
```

### Log Levels

- `debug`: Detailed debugging information, including all verification steps
- `info`: General information, including requests and responses
- `warn`: Warning information, including verification failures
- `error`: Error information, including system errors

## Error Codes

### Verification Errors

- `INVALID_X402_VERSION`: X402 protocol version not supported
- `INVALID_PAYLOAD`: Request payload format error
- `UNSUPPORTED_SCHEME`: Unsupported payment scheme
- `INVALID_NETWORK`: Unsupported network
- `INVALID_EXACT_EVM_PAYLOAD_SIGNATURE`: Signature verification failed
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALUE`: Invalid authorization amount
- `INVALID_EXACT_EVM_PAYLOAD_RECIPIENT_MISMATCH`: Payee address mismatch
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_AFTER`: Authorization not yet valid
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_BEFORE`: Authorization expired
- `INSUFFICIENT_FUNDS`: Insufficient user balance

### Settlement Errors

- `INVALID_TRANSACTION_STATE`: Blockchain transaction failed or rejected
- `UNEXPECTED_VERIFY_ERROR`: Unexpected error during verification
- `UNEXPECTED_SETTLE_ERROR`: Unexpected error during settlement
- `UNKNOWN`: Unknown error

## Security Considerations

1. **Private Key Management**:
   - Never commit private keys to version control systems
   - Store private keys using environment variables or key management services
   - Regularly rotate private keys

2. **Network Security**:
   - Use HTTPS in production environments
   - Configure appropriate CORS policies
   - Implement rate limiting

3. **Input Validation**:
   - All inputs are strictly validated
   - Use type-safe validators

4. **Error Handling**:
   - Avoid leaking sensitive information in error messages
   - Log detailed error information for debugging


## Useful links

* [Official documentation](https://x402.gitbook.io/x402)
* [X402 GitHub](https://github.com/coinbase/x402)
* [White paper](https://www.x402.org/x402-whitepaper.pdf)
* [Specifications](https://github.com/coinbase/x402/blob/main/specs/schemes/exact/scheme_exact_evm.md)
* [Examples](https://github.com/coinbase/x402/tree/main/examples/typescript)
* [CDP faucet](https://portal.cdp.coinbase.com/products/faucet)
* [Circle faucet](https://faucet.circle.com/)
* [Base Sepolia Testnet Explorer](https://sepolia.basescan.org/)