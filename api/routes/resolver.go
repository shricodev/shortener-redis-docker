package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"

	"github.com/shricodev/shortener-redis-docker/database"
)

// Resolve the URL from the shortened URL and redirect to the original URL
func ResolveURL(c *fiber.Ctx) error {
	url := c.Params("url")
	rdb := database.CreateClient(0)
	defer rdb.Close()

	val, err := rdb.Get(database.Ctx, url).Result()
	if err == redis.Nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
		// Not able to connect to the Redis server or some error occurred
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	rInr := database.CreateClient(1)
	defer rInr.Close()

	_ = rInr.Incr(database.Ctx, "counter")

	return c.Redirect(val, fiber.StatusMovedPermanently)
}
