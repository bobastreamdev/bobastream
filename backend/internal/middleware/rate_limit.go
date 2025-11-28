package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"bobastream/config"
)

// RateLimitAuth creates rate limiter for auth endpoints
func RateLimitAuth() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.GlobalConfig.RateLimit.Auth,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests, please try again later",
			})
		},
	})
}

// RateLimitAPI creates rate limiter for API endpoints
func RateLimitAPI() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.GlobalConfig.RateLimit.API,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			// Use user_id if authenticated, otherwise IP
			if userID := c.Locals("user_id"); userID != nil {
				return userID.(string)
			}
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "API rate limit exceeded",
			})
		},
	})
}

// RateLimitStream creates rate limiter for streaming endpoints
func RateLimitStream() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        config.GlobalConfig.RateLimit.Stream,
		Expiration: 1 * time.Hour,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Streaming rate limit exceeded",
			})
		},
	})
}