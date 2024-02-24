package routes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/shricodev/shortener-redis-docker/database"
	"github.com/shricodev/shortener-redis-docker/helpers"
)

const (
	RateLimitDuration    = 24 * time.Hour
	CustomShortURLLength = 3
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

	// Check if the user has reached the rate limit.
	r2 := database.CreateClient(1)
	defer r2.Close()

	val, err := r2.Get(database.Ctx, c.IP()).Result()
	// The user has not used the application in the last {RateLimitDuration} time
	if err != redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), RateLimitDuration).Err()
	} else {
		// The user has used the application in the last {RateLimitDuration} time
		valInt, _ := strconv.Atoi(val)

		if valInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()

			// Convert the limit duration into hours and minutes
			hours := int(limit.Hours())
			minutes := int(limit.Minutes()) % 60

			// Construct a readable string for the reset time
			resetTime := fmt.Sprintf("%d hours %d minutes", hours, minutes)

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error":    "Rate limit exceeded",
				"reset_in": resetTime,
			})
		}
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// If the user tries to use our own domain, we will return an error.
	if !helpers.CheckForApplicationMatch(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Invalid URL",
		})
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	var short_url string

	// Trim spaces and check if the customShortURL length is at least {CustomShortURLLength}
	trimmedCustomShortURL := strings.TrimSpace(body.CustomShortURL)
	if len(trimmedCustomShortURL) == 0 {
		trimmedCustomShortURL = uuid.New().String()[:6]
	} else if len(trimmedCustomShortURL) >= CustomShortURLLength {
		short_url = trimmedCustomShortURL
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Custom short URL must be at least 3 characters long and cannot contain spaces",
		})
	}

	// Check if the custom short URL is already in use
	r := database.CreateClient(0)
	defer r.Close()

	// Check if the custom short URL is already in use
	in_use_val, _ := r.Get(database.Ctx, short_url).Result()
	if strings.TrimSpace(in_use_val) != "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": "Custom short URL is already in use",
		})
	}

	// If the user has not set an expiration time, we will set it to 1 day.
	if body.Expiration == 0 {
		body.Expiration = 24 * time.Hour
	}

	// Save the URL to the database
	err = r.Set(database.Ctx, short_url, body.URL, body.Expiration).Err()
  if err != nil {
    return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
      "error": "Failed to save URL",
    })
  }

	r2.Decr(database.Ctx, c.IP())
	return nil
}
