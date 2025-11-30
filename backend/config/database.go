package config

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitDatabase initializes database connection with GORM
func InitDatabase() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta",
		GlobalConfig.Database.Host,
		GlobalConfig.Database.Port,
		GlobalConfig.Database.User,
		GlobalConfig.Database.Password,
		GlobalConfig.Database.Name,
		GlobalConfig.Database.SSLMode,
	)

	// Configure GORM logger
	var gormLogger logger.Interface
	if GlobalConfig.App.Env == "development" {
		gormLogger = logger.Default.LogMode(logger.Info)
	} else {
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// Open connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			// Use Jakarta timezone
			loc, _ := time.LoadLocation("Asia/Jakarta")
			return time.Now().In(loc)
		},
	})
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	// Get underlying SQL DB for connection pooling
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %w", err)
	}

	// ✅ PRODUCTION-READY CONNECTION POOL SETTINGS
	// Formula: MaxOpenConns = (Expected concurrent requests / Avg query time) * 1.5
	// For 1000 req/s with 100ms avg query: (1000 * 0.1) * 1.5 = 150 connections
	sqlDB.SetMaxOpenConns(100)                    // ✅ 100 concurrent connections (increased from 25)
	sqlDB.SetMaxIdleConns(25)                     // ✅ 25 idle connections (increased from 10)
	sqlDB.SetConnMaxLifetime(1 * time.Hour)       // ✅ Reuse connections longer (increased from 5 minutes)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)    // ✅ Close idle connections after 10 minutes

	// Test connection
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("✅ Database connected successfully")
	log.Printf("📊 Connection pool settings: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%s\n",
		100, 25, 1*time.Hour)
	
	return nil
}

// CloseDatabase closes database connection
func CloseDatabase() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}