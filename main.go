package main

import (
	"log"

	"github.com/WaffleHacks/mailer/logging"
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
}
