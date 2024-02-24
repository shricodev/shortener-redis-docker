package initializers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func LoadAndCheckEnv() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal().Err(err).Msg("Failed to load environment variables from .env file")
	}

	// Set default values for environment variables if they are not set
	setEnvIfNotSet("APPLICATION_PORT", "3000")
	setEnvIfNotSet("REDIS_PORT", "6379")
	setEnvIfNotSet("API_QUOTA", "10") // Assuming a default quota value, adjust as necessary
	setEnvIfNotSet("REDIS_PASSWORD", "")
	setEnvIfNotSet("APPLICATION_HOST", "localhost")
}

// setEnvIfNotSet sets the environment variable to the provided default value if it's not already set.
func setEnvIfNotSet(key, defaultValue string) {
	if _, exists := os.LookupEnv(key); !exists {
		os.Setenv(key, defaultValue)
	}
}
