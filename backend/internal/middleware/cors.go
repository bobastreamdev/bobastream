package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"bobastream/config"
)

// SetupCORS configures CORS middleware
func SetupCORS() fiber.Handler {
	origins := strings.Split(config.GlobalConfig.CORS.AllowedOrigins, ",")

	return cors.New(cors.Config{
		AllowOrigins:     strings.Join(origins, ","),
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
		MaxAge:           3600,
	})
}