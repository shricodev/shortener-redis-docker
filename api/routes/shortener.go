package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/shricodev/shortener-redis-docker/helpers"
)

type request struct {
	URL            string        `json:"url"`
	CustomShortURL string        `json:"custom_short"`
	Expiration     time.Duration `json:"expiration"`
}

type response struct {
	Url             string        `json:"url"`
	CustomShortURL  string        `json:"custom_short"`
	Expiration      time.Duration `json:"expiration"`
	XRateRemains    int           `json:"x-rate-remain"`
	XRateLimitReset time.Duration `json:"x-rate-limit-reset"`
}

func ShortenURL(c *fiber.Ctx) error {
	body := new(request)
	if !helpers.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	if !helpers.CheckForApplicationMatch(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)
	return nil
}
