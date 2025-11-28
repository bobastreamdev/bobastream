package config

import (
	"time"
)

// ParseDuration safely parses duration string
func GetAccessTokenDuration() time.Duration {
	duration, err := time.ParseDuration(GlobalConfig.JWT.AccessExpiry)
	if err != nil {
		return 15 * time.Minute // Default 15 minutes
	}
	return duration
}

func GetRefreshTokenDuration() time.Duration {
	duration, err := time.ParseDuration(GlobalConfig.JWT.RefreshExpiry)
	if err != nil {
		return 168 * time.Hour // Default 7 days
	}
	return duration
}