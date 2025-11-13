package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"x402-facilitator-go/internal/config"
	"x402-facilitator-go/internal/handlers"
	"x402-facilitator-go/internal/middleware"
	"x402-facilitator-go/internal/service"
	"x402-facilitator-go/internal/verifier"
	"x402-facilitator-go/internal/verifier/exact"
	"x402-facilitator-go/internal/web3"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	logger, err := initLogger(cfg.Logging)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	// Validate configuration
	if err := cfg.Validate(); err != nil {
		logger.Fatal("Invalid configuration", zap.Error(err))
	}

	logger.Info("Starting X402 Facilitator",
		zap.String("version", "1.0.0"),
		zap.String("address", cfg.Server.Address()),
	)

	// Initialize Web3 client from network configuration
	web3Client, err := web3.NewClient(cfg.Networks.NetworkInfos, logger)
	if err != nil {
		logger.Fatal("Failed to initialize Web3 client", zap.Error(err))
	}
	defer func() {
		if err := web3Client.Close(); err != nil {
			logger.Error("Error closing Web3 client", zap.Error(err))
		}
	}()

	// Initialize verifiers
	verifiers := []verifier.Verifier{
		exact.NewGlobalVerifier(logger),
		exact.NewPaymentContextVerifier(logger, web3Client),
		exact.NewEIP3009AssetVerifier(logger, web3Client),
		exact.NewSignatureVerifier(logger, web3Client),
		exact.NewUserBalanceVerifier(logger, web3Client),
	}

	// Initialize services
	verifyService := service.NewVerifyService(verifiers, logger)
	settleService := service.NewSettleService(verifyService, web3Client, cfg.X402.FacilitatorPrivateKey, logger)
	supportedService := service.NewSupportedService(cfg.Networks.NetworkInfos)

	// Initialize handlers
	verifyHandler := handlers.NewVerifyHandler(verifyService, logger)
	settleHandler := handlers.NewSettleHandler(settleService, logger)
	supportedHandler := handlers.NewSupportedHandler(supportedService, logger)

	// Setup router
	router := setupRouter(logger, verifyHandler, settleHandler, supportedHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	logger.Info("Server started successfully")

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exited")
}

// setupRouter configures the HTTP router
func setupRouter(
	logger *zap.Logger,
	verifyHandler *handlers.VerifyHandler,
	settleHandler *handlers.SettleHandler,
	supportedHandler *handlers.SupportedHandler,
) *gin.Engine {
	// Set Gin mode
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	// Apply middleware
	router.Use(middleware.Logger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// API routes
	api := router.Group("")
	{
		api.POST("/verify", verifyHandler.Verify)
		api.POST("/settle", settleHandler.Settle)
		api.GET("/supported", supportedHandler.Supported)
	}

	return router
}

// initLogger initializes the logger
func initLogger(loggingCfg config.LoggingConfig) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	if loggingCfg.Level != "" {
		if lvl, err := zapcore.ParseLevel(loggingCfg.Level); err == nil {
			cfg.Level = zap.NewAtomicLevelAt(lvl)
		}
	}

	switch loggingCfg.Format {
	case "console":
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
	default:
		cfg.Encoding = "json"
	}

	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	return logger, nil
}
