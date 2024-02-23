package initializers

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func CheckEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal().Msg("There was an error starting the server")
	}

	_, ok := os.LookupEnv("APPLICATION_PORT")

	// If the APPLICATION_PORT environment variable is not set, we will set it to 3000
	if !ok {
		os.Setenv("APPLICATION_PORT", "3000")
	}
}
