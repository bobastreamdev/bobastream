package main

import (
	"bobastream/config"
	"bobastream/internal/models"
	"bobastream/internal/utils"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Initialize database
	if err := config.InitDatabase(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer config.CloseDatabase()

	db := config.DB

	log.Println("🌱 Starting database seeding...")

	// Seed in order
	if err := seedCategories(db); err != nil {
		log.Fatal("Failed to seed categories:", err)
	}

	if err := seedAdminUser(db); err != nil {
		log.Fatal("Failed to seed admin user:", err)
	}

	if err := seedSampleAds(db); err != nil {
		log.Fatal("Failed to seed sample ads:", err)
	}

    if err := seedSamplePCloudCredential(db); err != nil {
        log.Fatal("Failed to seed pCloud credential:", err)
    }
    
    if err := seedSampleVideo(db); err != nil {
        log.Fatal("Failed to seed sample video:", err)
    }
    
    log.Println("✅ Database seeding completed successfully!")
}

// seedCategories seeds default categories
func seedCategories(db *gorm.DB) error {
	log.Println("📂 Seeding categories...")

	categories := []models.Category{
		{
			Name:         "Mahasiswi",
			Slug:         "mahasiswi",
			Description:  "Konten mahasiswi Indonesia",
			Icon:         "🎓",
			DisplayOrder: 1,
			IsActive:     true,
		},
		{
			Name:         "SMA",
			Slug:         "sma",
			Description:  "Konten anak SMA",
			Icon:         "📚",
			DisplayOrder: 2,
			IsActive:     true,
		},
		{
			Name:         "ABG",
			Slug:         "abg",
			Description:  "Anak baru gede",
			Icon:         "👧",
			DisplayOrder: 3,
			IsActive:     true,
		},
		{
			Name:         "Tante",
			Slug:         "tante",
			Description:  "Wanita dewasa",
			Icon:         "👩",
			DisplayOrder: 4,
			IsActive:     true,
		},
		{
			Name:         "Jilbab",
			Slug:         "jilbab",
			Description:  "Berjilbab",
			Icon:         "🧕",
			DisplayOrder: 5,
			IsActive:     true,
		},
		{
			Name:         "Indo",
			Slug:         "indo",
			Description:  "Indonesia asli",
			Icon:         "🇮🇩",
			DisplayOrder: 6,
			IsActive:     true,
		},
		{
			Name:         "Colmek",
			Slug:         "colmek",
			Description:  "Coli memek",
			Icon:         "💦",
			DisplayOrder: 7,
			IsActive:     true,
		},
		{
			Name:         "Live",
			Slug:         "live",
			Description:  "Live show",
			Icon:         "🔴",
			DisplayOrder: 8,
			IsActive:     true,
		},
		{
			Name:         "Viral",
			Slug:         "viral",
			Description:  "Video viral terbaru",
			Icon:         "🔥",
			DisplayOrder: 9,
			IsActive:     true,
		},
		{
			Name:         "Premium",
			Slug:         "premium",
			Description:  "Konten premium eksklusif",
			Icon:         "⭐",
			DisplayOrder: 10,
			IsActive:     true,
		},
	}

	for _, category := range categories {
		// Check if category already exists
		var existing models.Category
		err := db.Where("slug = ?", category.Slug).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Create new category
			if err := db.Create(&category).Error; err != nil {
				return fmt.Errorf("failed to create category %s: %w", category.Name, err)
			}
			log.Printf("  ✅ Created category: %s", category.Name)
		} else if err != nil {
			return fmt.Errorf("failed to check category %s: %w", category.Name, err)
		} else {
			log.Printf("  ⏭️  Category already exists: %s", category.Name)
		}
	}

	return nil
}

// seedAdminUser seeds default admin user
func seedAdminUser(db *gorm.DB) error {
	log.Println("👤 Seeding admin user...")

	// Default admin credentials
	email := "admin@bobastream.com"
	username := "admin"
	password := "admin123" // Change this in production!

	// Check if admin already exists
	var existing models.User
	err := db.Where("email = ?", email).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		// Hash password
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Create admin user
		admin := models.User{
			Email:        email,
			Username:     username,
			PasswordHash: hashedPassword,
			Role:         models.RoleAdmin,
		}

		if err := db.Create(&admin).Error; err != nil {
			return fmt.Errorf("failed to create admin user: %w", err)
		}

		log.Printf("  ✅ Created admin user:")
		log.Printf("     Email: %s", email)
		log.Printf("     Username: %s", username)
		log.Printf("     Password: %s", password)
		log.Printf("     ⚠️  CHANGE THIS PASSWORD IN PRODUCTION!")
	} else if err != nil {
		return fmt.Errorf("failed to check admin user: %w", err)
	} else {
		log.Printf("  ⏭️  Admin user already exists: %s", email)
	}

	return nil
}

// seedSampleAds seeds sample ads for testing
func seedSampleAds(db *gorm.DB) error {
	log.Println("📺 Seeding sample ads...")

	ads := []models.Ad{
		{
			Title:            "Sample Preroll Ad",
			AdType:           models.AdTypePreroll,
			ContentURL:       "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_1mb.mp4",
			RedirectURL:      "https://example.com",
			DurationSeconds:  7,
			DisplayFrequency: 1,
			Priority:         1,
			IsActive:         true,
		},
		{
			Title:            "Sample Banner Ad",
			AdType:           models.AdTypeBanner,
			ContentURL:       "https://via.placeholder.com/728x90?text=Banner+Ad",
			RedirectURL:      "https://example.com",
			DurationSeconds:  0,
			DisplayFrequency: 1,
			Priority:         1,
			IsActive:         true,
		},
		{
			Title:            "Sample Popup Ad",
			AdType:           models.AdTypePopup,
			ContentURL:       "https://via.placeholder.com/300x250?text=Popup+Ad",
			RedirectURL:      "https://example.com",
			DurationSeconds:  5,
			DisplayFrequency: 3,
			Priority:         1,
			IsActive:         false, // Disabled by default
		},
	}

	for _, ad := range ads {
		// Check if ad with same title exists
		var existing models.Ad
		err := db.Where("title = ?", ad.Title).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&ad).Error; err != nil {
				return fmt.Errorf("failed to create ad %s: %w", ad.Title, err)
			}
			log.Printf("  ✅ Created ad: %s (%s)", ad.Title, ad.AdType)
		} else if err != nil {
			return fmt.Errorf("failed to check ad %s: %w", ad.Title, err)
		} else {
			log.Printf("  ⏭️  Ad already exists: %s", ad.Title)
		}
	}

	return nil
}

// Optional: Seed sample pCloud credential
func seedSamplePCloudCredential(db *gorm.DB) error {
	log.Println("☁️  Seeding sample pCloud credential...")

	// Check if any credential exists
	var count int64
	if err := db.Model(&models.PCloudCredential{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("  ⏭️  pCloud credentials already exist")
		return nil
	}

	credential := models.PCloudCredential{
		AccountName:    "Sample Account",
		APIToken:       "YOUR_PCLOUD_API_TOKEN_HERE",
		StorageUsedGB:  0,
		StorageLimitGB: 500,
		IsActive:       false, // Disabled by default
	}

	if err := db.Create(&credential).Error; err != nil {
		return fmt.Errorf("failed to create pCloud credential: %w", err)
	}

	log.Println("  ✅ Created sample pCloud credential (INACTIVE)")
	log.Println("     ⚠️  Update API token and activate in admin panel")

	return nil
}

// Optional: Seed sample video for testing (requires pCloud credential)
func seedSampleVideo(db *gorm.DB) error {
	log.Println("🎥 Seeding sample video...")

	// Check if any pCloud credential exists
	var credential models.PCloudCredential
	err := db.Where("is_active = ?", true).First(&credential).Error
	if err != nil {
		log.Println("  ⏭️  No active pCloud credential found, skipping video seed")
		return nil
	}

	// Check if category exists
	var category models.Category
	err = db.Where("slug = ?", "viral").First(&category).Error
	if err != nil {
		log.Println("  ⏭️  Category 'viral' not found, skipping video seed")
		return nil
	}

	// Check if sample video exists
	var existing models.Video
	err = db.Where("title = ?", "Sample Test Video").First(&existing).Error
	if err != gorm.ErrRecordNotFound {
		log.Println("  ⏭️  Sample video already exists")
		return nil
	}

	now := time.Now()
	expiresAt := now.Add(24 * time.Hour)

	video := models.Video{
		Title:              "Sample Test Video",
		Description:        "This is a sample video for testing purposes",
		ThumbnailURL:       "https://via.placeholder.com/640x360?text=Sample+Video",
		SourceURL:          "https://sample-videos.com/video321/mp4/720/big_buck_bunny_720p_1mb.mp4",
		SourceURLExpiresAt: &expiresAt,
		DurationSeconds:    60,
		FileSizeMB:         1.5,
		PCloudFileID:       "sample123",
		PCloudCredentialID: credential.ID,
		CategoryID:         &category.ID,
		Tags:               []string{"sample", "test", "demo"},
		ViewCount:          0,
		LikeCount:          0,
		IsPublished:        true,
		PublishedAt:        &now,
	}

	if err := db.Create(&video).Error; err != nil {
		return fmt.Errorf("failed to create sample video: %w", err)
	}

	// Create wrapper link
	wrapperLink := models.WrapperLink{
		VideoID:      video.ID,
		WrapperToken: uuid.New().String(),
	}

	if err := db.Create(&wrapperLink).Error; err != nil {
		return fmt.Errorf("failed to create wrapper link: %w", err)
	}

	log.Printf("  ✅ Created sample video with wrapper token: %s", wrapperLink.WrapperToken)
	log.Printf("     Watch at: http://localhost:8080/watch/%s", wrapperLink.WrapperToken)

	return nil
}