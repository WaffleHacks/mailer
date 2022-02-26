package main

import (
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

// Config contains all the runtime configuration information
type Config struct {
	GRPCAddress string
	HTTPAddress string

	LogLevel    zap.AtomicLevel
	Development bool
}

// ReadConfig extracts all the configuration options from the environment variables
func ReadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	address := getEnvOrDefault("MAILER_ADDRESS", "127.0.0.1")
	grpcPort := getEnvOrDefault("MAILER_GRPC_PORT", "9000")
	httpPort := getEnvOrDefault("MAILER_HTTP_PORT", "8000")

	rawLevel := getEnvOrDefault("MAILER_LOG_LEVEL", "info")
	level := zap.NewAtomicLevel()
	if err := level.UnmarshalText([]byte(rawLevel)); err != nil {
		return nil, err
	}

	rawDevelopment := strings.ToLower(getEnvOrDefault("MAILER_DEVELOPMENT", "no"))
	development := rawDevelopment == "y" || rawDevelopment == "yes" || rawDevelopment == "t" || rawDevelopment == "true"

	return &Config{
		GRPCAddress: net.JoinHostPort(address, grpcPort),
		HTTPAddress: net.JoinHostPort(address, httpPort),
		LogLevel:    level,
		Development: development,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
