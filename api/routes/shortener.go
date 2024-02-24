package routes

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/shricodev/shortener-redis-docker/database"
	"github.com/shricodev/shortener-redis-docker/helpers"
)

const (
	rateLimitDuration   = 24 * time.Hour
	customShortUrlIdLen = 6
)

type request struct {
	URL              string        `json:"url"`
	CustomShortUrlId string        `json:"custom_short"`
	Expiration       time.Duration `json:"expiration"`
}

type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiration      time.Duration `json:"expiration"`
	XRateRemains    int           `json:"x-rate-remain"`
	XRateLimitReset time.Duration `json:"x-rate-limit-reset"`
}

// ShortenURL orchestrates the shortening process by validating the request, generating the short URL, and preparing the response.
func ShortenURL(c *fiber.Ctx) error {
	body := new(request)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON body",
		})
	}

	if err := ValidateRequest(c, body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	shortUrlID, err := GenerateShortURL(c, body)
	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	r2 := database.CreateClient(1)
	defer r2.Close()

	quotaLeft, _ := r2.Get(database.Ctx, c.IP()).Result()
	ttl, _ := r2.TTL(database.Ctx, c.IP()).Result()

	resp, err := PrepareResponse(c, body, shortUrlID, quotaLeft, ttl)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to prepare response",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ValidateRequest validates the request and checks for rate limiting.
func ValidateRequest(c *fiber.Ctx, body *request) error {
	// Check if the user has reached the rate limit.
	r2 := database.CreateClient(1)
	defer r2.Close()

	quotaLeft, err := r2.Get(database.Ctx, c.IP()).Result()
	if err == redis.Nil {
		_ = r2.Set(database.Ctx, c.IP(), os.Getenv("API_QUOTA"), rateLimitDuration).Err()
	} else {
		quotaLeftInt, _ := strconv.Atoi(quotaLeft)
		if quotaLeftInt <= 0 {
			limit, _ := r2.TTL(database.Ctx, c.IP()).Result()
			return fmt.Errorf("Rate limit exceeded. Reset in %v", limit)
		}
	}

	if !govalidator.IsURL(body.URL) {
		return fmt.Errorf("Invalid URL")
	}

	if !helpers.CheckForApplicationMatch(body.URL) {
		return fmt.Errorf("Provided URL is not allowed")
	}

	body.URL = helpers.EnforceHTTP(body.URL)

	return nil
}

// GenerateShortURL generates a short URL and saves it to the database.
func GenerateShortURL(c *fiber.Ctx, body *request) (string, error) {
	var shortUrlID string

	trimmedCustomShortUrl := strings.TrimSpace(body.CustomShortUrlId)
	if len(trimmedCustomShortUrl) == 0 {
		trimmedCustomShortUrl = uuid.New().String()[:6]
	} else if len(trimmedCustomShortUrl) >= customShortUrlIdLen {
		shortUrlID = trimmedCustomShortUrl
	} else {
		return "", fmt.Errorf("Custom short URL must be at least %v characters long", customShortUrlIdLen)
	}

	r := database.CreateClient(0)
	defer r.Close()

	inUseVal, _ := r.Get(database.Ctx, shortUrlID).Result()
	if strings.TrimSpace(inUseVal) != "" {
		return "", fmt.Errorf("Custom short URL is already in use")
	}

	if body.Expiration == 0 {
		body.Expiration = 24 * time.Hour
	}

	err := r.Set(database.Ctx, shortUrlID, body.URL, body.Expiration).Err()
	if err != nil {
		return "", fmt.Errorf("Failed to save URL")
	}

	return shortUrlID, nil
}

// PrepareResponse prepares the response to be sent back to the client.
func PrepareResponse(c *fiber.Ctx, body *request, shortUrlID string, quotaLeft string, ttl time.Duration) (*response, error) {
	resp := &response{
		URL:             body.URL,
		CustomShort:     "",
		Expiration:      body.Expiration,
		XRateRemains:    10,
		XRateLimitReset: 24,
	}

	quotaLeftInt, _ := strconv.Atoi(quotaLeft)
	resp.XRateRemains = quotaLeftInt
	resp.XRateLimitReset = ttl / time.Nanosecond / time.Minute / time.Hour

	applicationHost := os.Getenv("APPLICATION_HOST")
	applicationPort := os.Getenv("APPLICATION_PORT")

	if strings.TrimSpace(applicationPort) != "" {
		resp.CustomShort = fmt.Sprintf("%s:%s/%s", applicationHost, applicationPort, shortUrlID)
	} else {
		resp.CustomShort = fmt.Sprintf("%s/%s", applicationHost, shortUrlID)
	}

	return resp, nil
}
