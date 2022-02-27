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

	"github.com/WaffleHacks/mailer/daemon"
	"github.com/WaffleHacks/mailer/providers"
)

// Config contains all the runtime configuration information
type Config struct {
	GRPCAddress string
	HTTPAddress string

	LogLevel    zap.AtomicLevel
	Development bool

	Workers int

	Providers []*daemon.Matcher
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

	// Register all the providers
	enabledProviders := strings.Split(os.Getenv("MAILER_PROVIDERS"), ",")
	var configuredProviders []*daemon.Matcher
	for _, rawId := range enabledProviders {
		if len(rawId) == 0 {
			continue
		}
		id := strings.TrimSpace(rawId)
		envId := strings.ToUpper(id)

		// Create the provider
		typeName := os.Getenv(fmt.Sprintf("MAILER_PROVIDER_%s_TYPE", envId))
		provider, err := providers.Get(id, typeName)
		if err != nil {
			return nil, err
		} else if provider == nil {
			return nil, fmt.Errorf("unknown provider type %q for %q", typeName, id)
		}

		// Get the number of workers for the provider
		workers, err := strconv.Atoi(getEnvOrDefault(fmt.Sprintf("MAILER_PROVIDER_%s_WORKERS", envId), "1"))
		if err != nil {
			return nil, fmt.Errorf("failed to parse number of workers for provider %s: %v", id, err)
		}

		// Determine which providers should match
		pattern := getEnvOrDefault(fmt.Sprintf("MAILER_PROVIDER_%s_MATCHES", envId), "*")
		matcher, err := daemon.NewMatcher(id, workers, provider, pattern)
		if err != nil {
			return nil, fmt.Errorf("invalid patter for matching provider %q: %v", id, err)
		}

		configuredProviders = append(configuredProviders, matcher)
	}
	if len(configuredProviders) == 0 {
		return nil, errors.New("at least 1 provider must be configured")
	}

	return &Config{
		GRPCAddress: net.JoinHostPort(address, grpcPort),
		HTTPAddress: net.JoinHostPort(address, httpPort),
		LogLevel:    level,
		Development: development,
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
