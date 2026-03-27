package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
	"github.com/pilgrim/gcp-camunda8-loan-orchestrator/internal/config"
	"github.com/pilgrim/gcp-camunda8-loan-orchestrator/internal/handlers"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer func() { _ = logger.Sync() }()

	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("load config", zap.Error(err))
	}

	// OAuth uses ZEEBE_CLIENT_ID, ZEEBE_CLIENT_SECRET, ZEEBE_AUTHORIZATION_SERVER_URL from the environment
	// (set by config.Load and/or .env). Gateway address drives TLS and default token audience.
	client, err := zbc.NewClient(&zbc.ClientConfig{
		GatewayAddress: cfg.ZeebeAddress,
	})
	if err != nil {
		logger.Fatal("create zeebe client", zap.Error(err))
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("close zeebe client", zap.Error(closeErr))
		}
	}()

	enrichment := handlers.NewEnrichmentHandler(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	jobWorker := client.NewJobWorker().
		JobType(cfg.EnrichmentJobType).
		Handler(enrichment.Handle).
		Name("credit-card-enrichment-worker").
		Open()

	logger.Info("worker started",
		zap.String("job_type", cfg.EnrichmentJobType),
		zap.String("gateway", cfg.ZeebeAddress),
	)

	<-ctx.Done()
	logger.Info("shutdown signal received, stopping job worker")

	jobWorker.Close()
	jobWorker.AwaitClose()
	logger.Info("worker stopped")
}
