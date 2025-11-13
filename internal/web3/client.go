package web3

import (
	"fmt"
	"math/big"
	"sync"
	"x402-facilitator-go/internal/config"

	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"
)

// Client manages Web3 clients for different networks
type Client struct {
	ClientInfo map[string]ClientInfo
	logger     *zap.Logger
	mu         sync.RWMutex
}

type ClientInfo struct {
	client  *ethclient.Client
	rpcURL  string
	chainID *big.Int
}

// NewClient creates a new Web3 client manager
func NewClient(networkInfo []config.NetworkInfo, logger *zap.Logger) (*Client, error) {
	clientMap := make(map[string]ClientInfo)

	for _, netInfo := range networkInfo {
		ethClient, err := ethclient.Dial(netInfo.RPCURL)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to %s at %s: %w", netInfo.Name, netInfo.RPCURL, err)
		}

		clientMap[netInfo.Name] = ClientInfo{
			client:  ethClient,
			rpcURL:  netInfo.RPCURL,
			chainID: big.NewInt(netInfo.ChainID),
		}
	}

	return &Client{
		ClientInfo: clientMap,
		logger:     logger,
	}, nil
}

// GetClient returns the Ethereum client for the specified network
func (c *Client) GetClient(networkName string) (*ethclient.Client, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clientInfo, ok := c.ClientInfo[string(networkName)]
	if !ok {
		return nil, fmt.Errorf("network %s not configured", networkName)
	}
	return clientInfo.client, nil
}

// Close closes all Web3 clients
func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for network, clientInfo := range c.ClientInfo {
		clientInfo.client.Close()
		c.logger.Info("Closed connection to network", zap.String("network", network))
	}

	return nil
}

// GetChainID returns the chain ID for the given network string
// This function is deprecated, use Client.GetChainID instead
func GetChainID(network string) (*big.Int, error) {
	return nil, fmt.Errorf("GetChainID is deprecated, use Client.GetChainID instead")
}

// GetChainID returns the chain ID for the specified network
func (c *Client) GetChainID(networkName string) (*big.Int, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	clientInfo, ok := c.ClientInfo[networkName]
	if !ok {
		return nil, fmt.Errorf("chain ID not configured for network %s", networkName)
	}

	return clientInfo.chainID, nil
}
