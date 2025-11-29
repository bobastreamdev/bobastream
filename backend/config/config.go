package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	PCloud    PCloudConfig
	CORS      CORSConfig
	RateLimit RateLimitConfig
	Cron      CronConfig
}

type AppConfig struct {
	Env  string
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	RefreshSecret string
	AccessExpiry  string
	RefreshExpiry string
}

type PCloudConfig struct {
	BaseURL string
}

type CORSConfig struct {
	AllowedOrigins string
}

type RateLimitConfig struct {
	Auth   int
	API    int
	Stream int
}

type CronConfig struct {
	RefreshLinks   string
	AggregateStats string
}

var GlobalConfig *Config

// LoadConfig loads environment variables and initializes config
func LoadConfig() error {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	GlobalConfig = &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
			Host: getEnv("APP_HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5555"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "boba_stream"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", ""),
			RefreshSecret: getEnv("JWT_REFRESH_SECRET", ""),
			AccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "15m"),
			RefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
		},
		PCloud: PCloudConfig{
			BaseURL: getEnv("PCLOUD_API_BASE_URL", "https://api.pcloud.com"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:5173"),
		},
		RateLimit: RateLimitConfig{
			Auth:   getEnvAsInt("RATE_LIMIT_AUTH", 5),
			API:    getEnvAsInt("RATE_LIMIT_API", 60),
			Stream: getEnvAsInt("RATE_LIMIT_STREAM", 100),
		},
		Cron: CronConfig{
			RefreshLinks:   getEnv("CRON_REFRESH_LINKS", "0 */6 * * *"),
			AggregateStats: getEnv("CRON_AGGREGATE_STATS", "0 0 * * *"),
		},
	}

	// Validate required configs
	if GlobalConfig.JWT.Secret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if GlobalConfig.JWT.RefreshSecret == "" {
		log.Fatal("JWT_REFRESH_SECRET is required")
	}

	return nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return fallback
}