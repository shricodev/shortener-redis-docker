package main

import (
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/rs/zerolog/log"

	"github.com/shricodev/shortener-redis-docker/initializers"
	"github.com/shricodev/shortener-redis-docker/routes"
)

func init() {
	initializers.CheckEnvVariables()
}

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL)
}

func main() {
	app := fiber.New()

	app.Use(logger.New())

	setupRoutes(app)

	if err := app.Listen(os.Getenv("APPLICATION_PORT")); err != nil {
		log.Fatal().Err(err).Msg("Failed to start application")
	}
}
