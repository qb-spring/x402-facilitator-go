package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/jinzhu/configor"
)

// ServerConfig holds server configuration
type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// X402Config holds X402 facilitator configuration
type X402Config struct {
	// FacilitatorPrivateKey is loaded from environment variable X402_FACILITATOR_PRIVATE_KEY
	// It is not read from YAML for security reasons
	FacilitatorPrivateKey string
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
}

// Config represents the YAML structure for unmarshaling
type Config struct {
	Server   ServerConfig  `yaml:"server"`
	X402     X402Config    `yaml:"x402"`
	Logging  LoggingConfig `yaml:"logging"`
	Networks NetworkConfig `yaml:"networks"`
}

// NetworkConfig is used for unmarshaling networks with string chainId
type NetworkConfig struct {
	NetworkInfos []NetworkInfo `yaml:"networkInfos"`
}

// NetworkInfo is used for unmarshaling chainId as string
type NetworkInfo struct {
	Name        string `yaml:"name"`
	RPCURL      string `yaml:"rpcURL"`
	ChainID     int64  `yaml:"chainId"`
	X402Version int16  `yaml:"X402Version"`
	Scheme      string `yaml:"scheme"`
}

// Load loads configuration from config.yaml file and environment variables
func Load() (*Config, error) {
	configPath := getConfigPath()

	// Load YAML config into intermediate structure
	config := &Config{}
	if err := configor.Load(config, configPath); err != nil {
		return nil, fmt.Errorf("failed to load config file %s: %w", configPath, err)
	}

	// Load private key from environment variable (security best practice)
	if privateKey := os.Getenv("X402_FACILITATOR_PRIVATE_KEY"); privateKey != "" {
		config.X402.FacilitatorPrivateKey = privateKey
	}

	return config, nil
}

// getConfigPath returns the path to the config file
// It tries multiple locations in order:
// 1. CONFIG_PATH environment variable
// 2. config.yaml in current working directory
// 3. config.yaml in project root (two levels up from internal/config)
func getConfigPath() string {
	// First, check environment variable
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	// Try current working directory
	if wd, err := os.Getwd(); err == nil {
		path := filepath.Join(wd, "config.yaml")
		if absPath, err := filepath.Abs(path); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}

	// Try project root (assuming we're in internal/config)
	// Get the directory of this source file
	if _, filename, _, ok := runtime.Caller(0); ok {
		// internal/config/config.go -> internal/config -> internal -> project root
		configDir := filepath.Dir(filename)
		projectRoot := filepath.Join(configDir, "..", "..")
		if absRoot, err := filepath.Abs(projectRoot); err == nil {
			path := filepath.Join(absRoot, "config.yaml")
			if _, err := os.Stat(path); err == nil {
				return path
			}
		}
	}

	// Fallback to relative path (current working directory)
	return "config.yaml"
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if c.X402.FacilitatorPrivateKey == "" {
		return fmt.Errorf("X402_FACILITATOR_PRIVATE_KEY environment variable is required")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	return nil
}

// Address returns the server address in the format "host:port"
func (s *ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}
