package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/WaffleHacks/mailer/daemon/providers"
)

// Config contains all the runtime configuration information
type Config struct {
	GRPCAddress string
	HTTPAddress string

	LogLevel    zap.AtomicLevel
	Development bool

	Workers int

	Providers map[string]providers.Provider
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

	workers, err := strconv.Atoi(getEnvOrDefault("MAILER_WORKERS", "1"))
	if err != nil {
		return nil, err
	}

	// Register all the providers
	enabledProviders := strings.Split(os.Getenv("MAILER_PROVIDERS"), ",")
	configuredProviders := make(map[string]providers.Provider)
	for _, rawId := range enabledProviders {
		if len(rawId) == 0 {
			continue
		}
		id := strings.TrimSpace(rawId)

		// Create the provider
		typeName := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_TYPE", strings.ToUpper(id)))
		provider, err := providers.Get(id, typeName)
		if err != nil {
			return nil, err
		} else if provider == nil {
			return nil, fmt.Errorf("unknown provider type %q for %q", typeName, id)
		}

		configuredProviders[id] = provider
	}
	if len(configuredProviders) == 0 {
		return nil, errors.New("at least 1 provider must be configured")
	}

	return &Config{
		GRPCAddress: net.JoinHostPort(address, grpcPort),
		HTTPAddress: net.JoinHostPort(address, httpPort),
		LogLevel:    level,
		Development: development,
		Workers:     workers,
		Providers:   configuredProviders,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}
