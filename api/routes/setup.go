package routes

import "github.com/gofiber/fiber/v2"

// It defines two routes: one for resolving a URL and another for shortening a URL
func SetupRoutes(app *fiber.App) {
	app.Get("/:url", ResolveURL)
	app.Post("/api/v1", ShortenURL)
}
