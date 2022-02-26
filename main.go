package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/WaffleHacks/mailer/logging"
	"github.com/WaffleHacks/mailer/rpc"
)

func main() {
	config, err := ReadConfig()
	if err != nil {
		log.Fatalf("failed to read configuration: %v\n", err)
	}

	logger, err := logging.New(config.LogLevel, config.Development)
	if err != nil {
		log.Fatalf("failed to initialize logging: %v\n", err)
	}
	defer logger.Sync()

	// Acquire the gRPC listener
	listener, err := net.Listen("tcp", config.GRPCAddress)
	if err != nil {
		logger.Fatal("failed to listen on address", zap.Error(err), zap.String("address", config.GRPCAddress))
	}

	// Start the gRPC server
	server := rpc.New()
	go func() {
		logger.Named("grpc").Info("listening and ready to handle requests", zap.String("address", config.GRPCAddress))
		if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
			logger.Named("grpc").Fatal("an error occurred while running the server", zap.Error(err))
		}
	}()

	// Start the HTTP gateway
	gateway, err := rpc.NewGateway(config.GRPCAddress, config.HTTPAddress)
	if err != nil {
		logger.Named("http").Fatal("failed to create HTTP gateway", zap.String("grpcAddress", config.GRPCAddress), zap.Error(err))
	}
	go func() {
		logger.Named("http").Info("listening and ready to handle requests", zap.String("address", config.HTTPAddress))
		if err := gateway.ListenAndServe(); err != nil && err != http.ErrServerClosed {
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

	server.GracefulStop()
	if err := gateway.Shutdown(shutdownCtx); err != nil {
		logger.Named("http").Fatal("failed to shutdown gateway", zap.Error(err))
	}

	logger.Info("shutdown complete. goodbye!")
}
