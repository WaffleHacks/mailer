package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/http"
	"github.com/WaffleHacks/mailer/tracing"

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/logging"
)

func main() {
	config, err := ReadConfig()
	if err != nil {
		log.Fatalf("failed to read configuration: %v\n", err)
	}

	// Initialize sentry if present
	if config.SentryDsn != nil {
		options := sentry.ClientOptions{
			Dsn:              config.SentryDsn.String(),
			TracesSampleRate: 0.5,
		}
		if config.Development {
			options.Environment = "development"
		} else {
			options.Environment = "production"
		}

		if err := sentry.Init(options); err != nil {
			log.Fatalf("failed to initialize sentry: %v\n", err)
		}
		defer sentry.Flush(time.Second * 3)
	}

	// Initialize tracing
	ctx := context.Background()
	provider, err := tracing.Initialize(ctx, config.Tracing, config.Development)
	if err != nil {
		log.Fatalf("failed to initalize tracing: %v\n", err)
	} else if provider != nil {
		defer func() {
			_ = provider.Shutdown(ctx)
		}()
	}

	logger, err := logging.New(config.LogLevel, config.Development)
	if err != nil {
		log.Fatalf("failed to initialize logging: %v\n", err)
	}
	defer logger.Sync()

	// Setup the mailer daemon
	mailer := daemon.New(ctx, config.Providers)

	// Start the HTTP server
	server := http.New(config.Address, mailer.Queue)
	if err != nil {
		logger.Named("http").Fatal("failed to create HTTP gateway", zap.String("address", config.Address), zap.Error(err))
	}
	go func() {
		logger.Named("http").Info("listening and ready to handle requests", zap.String("address", config.Address))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Named("http").Fatal("an error occurred while running the server", zap.Error(err))
		}
	}()

	// Register signal handlers for shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-shutdown

	// Set a 15s timeout for the shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Named("http").Fatal("failed to shutdown gateway", zap.Error(err))
	}
	if err := mailer.Shutdown(shutdownCtx); err != nil {
		logger.Named("daemon").Fatal("failed to shutdown daemon", zap.Error(err))
	}

	logger.Info("shutdown complete. goodbye!")
}
