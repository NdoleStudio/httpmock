package di

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

// Configuration is a struct that holds the configuration for the application.
type Configuration struct {
	UseOpenTelemetryLogger bool `env:"USE_OPEN_TELEMETRY_LOGGER"`
}

// LoadEnv will read your .env file(s) and load them into ENV for this process.
func LoadEnv(filenames ...string) {
	err := godotenv.Load(filenames...)
	if err != nil {
		log.Fatalf("Fatal: cannot load .env file: %v", err)
	}
}
