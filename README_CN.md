# X402 Facilitator

一个基于 Go 语言实现的 X402 支付协议 Facilitator（促进者）服务，支持基于 EIP-3009 标准的授权转账支付验证和结算。

## 目录

- [项目结构](#项目结构)
- [关于 X402](#关于-x402)
- [功能特性](#功能特性)
- [架构设计](#架构设计)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [开发指南](#开发指南)
- [错误码说明](#错误码说明)
- [安全注意事项](#安全注意事项)
- [实用链接](#实用链接)

## 项目结构

```
x402-facilitator/
├── cmd/
│   └── server/
│       └── main.go                    # 应用入口，初始化服务和路由
│
├── internal/
│   ├── config/
│   │   └── config.go                  # 配置管理，加载 YAML 配置和环境变量
│   │
│   ├── handlers/
│   │   ├── verify_handler.go          # 验证请求处理器 (POST /verify)
│   │   ├── settle_handler.go          # 结算请求处理器 (POST /settle)
│   │   └── supported_handler.go       # 支持查询处理器 (GET /supported)
│   │
│   ├── middleware/
│   │   ├── cors.go                    # CORS 跨域中间件
│   │   ├── logger.go                  # 请求日志中间件
│   │   └── recovery.go                # 错误恢复中间件
│   │
│   ├── models/
│   │   └── models.go                  # 数据模型定义（请求/响应结构体）
│   │
│   ├── service/
│   │   ├── verify_service.go          # 验证服务，协调多个验证器执行
│   │   ├── settle_service.go          # 结算服务，执行链上代币转账
│   │   └── supported_service.go       # 支持查询服务，返回支持的网络和方案
│   │
│   ├── util/
│   │   ├── eip3009/
│   │   │   └── eip3009.go             # EIP-3009 工具函数，计算授权哈希
│   │   └── eip712/
│   │       └── eip712.go              # EIP-712 工具函数，签名验证
│   │
│   ├── verifier/
│   │   ├── verifier.go                # 验证器接口定义
│   │   └── exact/
│   │       ├── global_verifier.go              # 全局验证器 (Order: 1)
│   │       ├── payment_context_verifier.go     # 支付上下文验证器 (Order: 2)
│   │       ├── eip3009_asset_verifier.go       # EIP-3009 资产验证器 (Order: 3)
│   │       ├── signature_verifier.go           # 签名验证器 (Order: 4)
│   │       └── user_balance_verifier.go       # 用户余额验证器 (Order: 5)
│   │
│   └── web3/
│       ├── client.go                  # Web3 客户端管理，支持多网络
│       └── contract/
│           └── EIP3009Token.go        # EIP-3009 合约 ABI 绑定
│
├── pkg/
│   └── errors/
│       └── errors.go                  # X402 错误码定义
│
├── config.yaml                         # 配置文件示例
├── go.mod                              # Go 模块定义
├── go.sum                              # 依赖校验和
└── README.md                           # 项目文档
```

### 目录说明

#### `cmd/server/`
应用入口点，负责：
- 加载配置
- 初始化日志系统
- 创建 Web3 客户端
- 注册验证器
- 初始化服务和处理器
- 设置 HTTP 路由
- 优雅关闭

#### `internal/config/`
配置管理模块：
- 从 YAML 文件加载配置
- 从环境变量加载敏感信息（如私钥）
- 配置验证

#### `internal/handlers/`
HTTP 请求处理器：
- `VerifyHandler`: 处理支付验证请求
- `SettleHandler`: 处理支付结算请求
- `SupportedHandler`: 返回支持的网络和方案列表

#### `internal/middleware/`
HTTP 中间件：
- `CORS`: 处理跨域请求
- `Logger`: 记录请求日志
- `Recovery`: 捕获 panic 并返回错误响应

#### `internal/models/`
数据模型定义：
- `VerifyRequest`: 验证请求结构
- `SettleRequest`: 结算请求结构
- `PaymentPayload`: 支付负载
- `PaymentRequirements`: 支付要求
- `Authorization`: 授权信息
- 响应结构体

#### `internal/service/`
业务逻辑层：
- `VerifyService`: 协调多个验证器按顺序执行验证
- `SettleService`: 执行链上代币转账
- `SupportedService`: 返回支持的网络配置

#### `internal/verifier/`
验证器模块，实现链式验证：
- `Verifier` 接口：定义验证器标准接口
- `exact/`: 实现 "exact" 支付方案的验证器
  - 按 `Order()` 方法定义的顺序执行
  - 任何验证器失败都会立即返回

#### `internal/util/`
工具函数：
- `eip3009/`: EIP-3009 标准相关工具
- `eip712/`: EIP-712 结构化数据签名工具

#### `internal/web3/`
区块链交互层：
- `Client`: 管理多个网络的以太坊客户端
- `contract/`: 智能合约 ABI 绑定

#### `pkg/errors/`
错误码定义：
- 定义所有 X402 协议错误码
- 提供错误码字符串转换方法

## 关于 X402

X402 是一个去中心化支付协议，旨在为 Web3 应用提供标准化的支付解决方案。该协议允许用户通过签名授权的方式完成支付，而无需在每次支付时手动确认交易。

### X402 核心概念

1. **支付授权（Payment Authorization）**：用户通过 EIP-712 签名创建支付授权，包含：
   - 付款人地址（From）
   - 收款人地址（To）
   - 支付金额（Value）
   - 有效期（ValidAfter, ValidBefore）
   - 随机数（Nonce）

2. **支付方案（Payment Scheme）**：当前支持 `exact` 方案，要求支付金额精确匹配授权金额。

3. **Facilitator（促进者）**：负责验证支付授权的有效性，并在验证通过后执行链上结算。

4. **EIP-3009 标准**：基于 ERC-20 的扩展标准，支持通过授权签名进行代币转账，无需用户每次手动确认。

## 功能特性

### 1. 多网络支持

支持多个 EVM 兼容网络，当前配置包括：
- Base Sepolia（测试网）
- Base Mainnet（主网）

可通过配置文件轻松添加更多网络。

### 2. 多层级验证

实现了完整的验证链，按顺序执行：

1. **全局验证（Global Verifier）**：验证请求格式和必填字段
2. **支付上下文验证（Payment Context Verifier）**：验证协议版本、方案、网络匹配性
3. **EIP-3009 资产验证（EIP-3009 Asset Verifier）**：验证代币合约是否支持 EIP-3009
4. **签名验证（Signature Verifier）**：使用 EIP-712 验证支付授权签名
5. **用户余额验证（User Balance Verifier）**：验证用户账户余额是否充足

### 3. 安全特性

- EIP-712 结构化数据签名验证
- 私钥通过环境变量管理，不存储在配置文件中
- 完整的错误处理和日志记录
- CORS 支持

### 4. 高可用性

- 优雅关闭（Graceful Shutdown）
- 上下文取消支持
- 结构化日志（JSON/Console 格式）
- 健康检查端点

## 架构设计

### 分层架构

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
│  (多个验证器按顺序执行)              │
└─────────────────────────────────────┘
              │
┌─────────────────────────────────────┐
│         Web3 Client Layer           │
│  (多网络 RPC 客户端管理)             │
└─────────────────────────────────────┘
```

### 验证器执行顺序

验证器按照 `Order()` 方法返回的顺序执行：

1. **Order 1**: `GlobalVerifier` - 全局格式验证
2. **Order 2**: `PaymentContextVerifier` - 支付上下文验证
3. **Order 3**: `EIP3009AssetVerifier` - 资产合约验证
4. **Order 4**: `SignatureVerifier` - 签名验证
5. **Order 5**: `UserBalanceVerifier` - 余额验证

任何验证器失败都会立即返回，不会继续执行后续验证。

### 数据流

```
客户端请求
    ↓
HTTP Handler (解析请求)
    ↓
Service Layer (业务逻辑)
    ↓
Verifier Chain (链式验证)
    ↓
Web3 Client (区块链交互)
    ↓
返回响应
```

## 快速开始

### 前置要求

- Go 1.21 或更高版本
- 访问 EVM 兼容网络的 RPC 节点
- Facilitator 私钥（用于执行结算交易）

### 安装

```bash
# 克隆仓库
git clone <repository-url>
cd x402-facilitator

# 安装依赖
go mod download
```

### 配置

1. 复制并编辑配置文件：

```bash
cp config.yaml config.yaml.local
```

2. 编辑 `config.yaml.local`，配置服务器和网络信息。

3. 设置环境变量：

```bash
export X402_FACILITATOR_PRIVATE_KEY="your_private_key_here"
```

### 运行

```bash
# 使用默认配置
go run cmd/server/main.go

# 或指定配置文件路径
CONFIG_PATH=./config.yaml.local go run cmd/server/main.go
```

### 构建

```bash
# 构建二进制文件
go build -o bin/x402-facilitator cmd/server/main.go

# 运行
./bin/x402-facilitator
```

## 配置说明

### 配置文件结构

```yaml
server:
  host: "0.0.0.0"      # 服务器监听地址
  port: 8081           # 服务器端口

logging:
  level: "info"        # 日志级别: debug, info, warn, error
  format: "json"       # 日志格式: json, console

networks:
  networkInfos:
    - name: "base-sepolia"           # 网络名称（用于 API 请求）
      rpcURL: "https://sepolia.base.org"  # RPC 节点 URL
      chainId: 84532                 # 链 ID
      X402Version: 1                 # 支持的 X402 协议版本
      scheme: "exact"                # 支持的支付方案
```

### 环境变量

- `X402_FACILITATOR_PRIVATE_KEY`：Facilitator 私钥（必需）
- `CONFIG_PATH`：配置文件路径（可选）

### 配置文件查找顺序

1. `CONFIG_PATH` 环境变量指定的路径
2. 当前工作目录的 `config.yaml`
3. 项目根目录的 `config.yaml`

## 开发指南

### 添加新的验证器

1. 在 `internal/verifier/exact/` 目录下创建新的验证器文件
2. 实现 `verifier.Verifier` 接口：
   ```go
   type Verifier interface {
       Verify(ctx context.Context, request *models.VerifyRequest) VerificationResult
       Type() VerificationStep
       Order() int
   }
   ```
3. 在 `cmd/server/main.go` 中注册验证器：
   ```go
   verifiers := []verifier.Verifier{
       // ... 现有验证器
       exact.NewYourVerifier(logger, web3Client),
   }
   ```

### 添加新的网络

在 `config.yaml` 中添加网络配置：

```yaml
networks:
  networkInfos:
    - name: "your-network"
      rpcURL: "https://your-rpc-url"
      chainId: 12345
      X402Version: 1
      scheme: "exact"
```

### 日志级别

- `debug`: 详细的调试信息，包括所有验证步骤
- `info`: 一般信息，包括请求和响应
- `warn`: 警告信息，包括验证失败
- `error`: 错误信息，包括系统错误


## 错误码说明

### 验证错误

- `INVALID_X402_VERSION`: X402 协议版本不支持
- `INVALID_PAYLOAD`: 请求负载格式错误
- `UNSUPPORTED_SCHEME`: 不支持的支付方案
- `INVALID_NETWORK`: 不支持的网络
- `INVALID_EXACT_EVM_PAYLOAD_SIGNATURE`: 签名验证失败
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALUE`: 授权金额无效
- `INVALID_EXACT_EVM_PAYLOAD_RECIPIENT_MISMATCH`: 收款人地址不匹配
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_AFTER`: 授权尚未生效
- `INVALID_EXACT_EVM_PAYLOAD_AUTHORIZATION_VALID_BEFORE`: 授权已过期
- `INSUFFICIENT_FUNDS`: 用户余额不足

### 结算错误

- `INVALID_TRANSACTION_STATE`: 区块链交易失败或被拒绝
- `UNEXPECTED_VERIFY_ERROR`: 验证过程中发生意外错误
- `UNEXPECTED_SETTLE_ERROR`: 结算过程中发生意外错误
- `UNKNOWN`: 未知错误

## 安全注意事项

1. **私钥管理**：
   - 永远不要将私钥提交到版本控制系统
   - 使用环境变量或密钥管理服务存储私钥
   - 定期轮换私钥

2. **网络安全**：
   - 在生产环境中使用 HTTPS
   - 配置适当的 CORS 策略
   - 实施速率限制

3. **输入验证**：
   - 所有输入都经过严格验证
   - 使用类型安全的验证器

4. **错误处理**：
   - 避免在错误消息中泄露敏感信息
   - 记录详细的错误日志用于调试

## 实用链接

* [Official documentation](https://x402.gitbook.io/x402)
* [X402 GitHub](https://github.com/coinbase/x402)
* [White paper](https://www.x402.org/x402-whitepaper.pdf)
* [Specifications](https://github.com/coinbase/x402/blob/main/specs/schemes/exact/scheme_exact_evm.md)
* [Examples](https://github.com/coinbase/x402/tree/main/examples/typescript)
* [CDP faucet](https://portal.cdp.coinbase.com/products/faucet)
* [Circle faucet](https://faucet.circle.com/)
* [Base Sepolia Testnet Explorer](https://sepolia.basescan.org/)