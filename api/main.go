package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/rs/zerolog/log"

	"github.com/shricodev/shortener-redis-docker/initializers"
	"github.com/shricodev/shortener-redis-docker/routes"
)

// Initialize the application by loading and checking environment variables
func init() {
	initializers.LoadAndCheckEnv()
}

func main() {
	// Create a new Fiber instance
	app := fiber.New()
	// Add logger middleware to log HTTP requests
	app.Use(logger.New())

	// Setup application routes
	routes.SetupRoutes(app)

	// Start the application on the port specified by the APPLICATION_PORT environment variable
	if err := app.Listen(os.Getenv("APPLICATION_PORT")); err != nil {
		// Log and exit the application if listening fails
		log.Fatal().Err(err).Msg("Failed to start application")
	}
}
