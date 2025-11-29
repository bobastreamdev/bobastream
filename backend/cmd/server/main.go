package main

import (
	"bobastream/config"
	"bobastream/internal/cron"
	"bobastream/internal/handlers"
	"bobastream/internal/middleware"
	"bobastream/internal/repositories"
	"bobastream/internal/services"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	cronpkg "github.com/robfig/cron/v3"
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

	// Initialize repositories
	userRepo := repositories.NewUserRepository(config.DB)
	videoRepo := repositories.NewVideoRepository(config.DB)
	wrapperRepo := repositories.NewWrapperLinkRepository(config.DB)
	videoViewRepo := repositories.NewVideoViewRepository(config.DB)
	videoLikeRepo := repositories.NewVideoLikeRepository(config.DB)
	categoryRepo := repositories.NewCategoryRepository(config.DB)
	adRepo := repositories.NewAdRepository(config.DB)
	adImpressionRepo := repositories.NewAdImpressionRepository(config.DB)
	analyticsRepo := repositories.NewAnalyticsRepository(config.DB)
	pcloudRepo := repositories.NewPCloudCredentialRepository(config.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	videoService := services.NewVideoService(videoRepo, wrapperRepo, videoViewRepo, videoLikeRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	adService := services.NewAdService(adRepo, adImpressionRepo)
	analyticsService := services.NewAnalyticsService(analyticsRepo, videoViewRepo, adImpressionRepo)
	pcloudService := services.NewPCloudService(pcloudRepo, videoRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	videoHandler := handlers.NewVideoHandler(videoService, categoryService)
	likeHandler := handlers.NewLikeHandler(videoService)
	streamHandler := handlers.NewStreamHandler(videoService)
	adHandler := handlers.NewAdHandler(adService)
	adminVideoHandler := handlers.NewAdminVideoHandler(videoService, pcloudService, categoryService)
	adminAdHandler := handlers.NewAdminAdHandler(adService)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService, videoService)
	pcloudHandler := handlers.NewPCloudHandler(pcloudService)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "BOBA STREAM API",
		ServerHeader: "Fiber",
		BodyLimit:    500 * 1024 * 1024, // 500MB for video uploads
	})

	// Global middleware
	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
	}))
	app.Use(middleware.SetupCORS())

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
			"app":    "BOBA STREAM API",
		})
	})

	// API routes
	api := app.Group("/api")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", middleware.RateLimitAuth(), authHandler.Register)
	auth.Post("/login", middleware.RateLimitAuth(), authHandler.Login)
	auth.Post("/logout", authHandler.Logout)
	auth.Post("/refresh", authHandler.RefreshToken)
	auth.Get("/me", middleware.AuthRequired(), authHandler.GetMe)

	// Public video routes
	videos := api.Group("/videos")
	videos.Get("/", videoHandler.GetFeed)
	videos.Get("/search", videoHandler.SearchVideos)
	videos.Get("/categories", videoHandler.GetCategories)
	videos.Get("/category/:categoryId", videoHandler.GetVideosByCategory)
	videos.Get("/:id", videoHandler.GetVideoByID)
	videos.Get("/:id/related", videoHandler.GetRelatedVideos)
	videos.Post("/:id/view", videoHandler.TrackView)

	// Like routes (auth required)
	videos.Post("/:id/like", middleware.AuthRequired(), likeHandler.LikeVideo)
	videos.Delete("/:id/like", middleware.AuthRequired(), likeHandler.UnlikeVideo)
	videos.Get("/:id/liked", middleware.AuthRequired(), likeHandler.CheckLiked)

	// User liked videos
	users := api.Group("/users", middleware.AuthRequired())
	users.Get("/me/likes", likeHandler.GetUserLikedVideos)

	// Ads routes (public)
	ads := api.Group("/ads")
	ads.Get("/preroll", adHandler.GetPrerollAd)
	ads.Get("/banner", adHandler.GetBannerAd)
	ads.Get("/popup", adHandler.GetPopupAd)
	ads.Post("/:id/impression", adHandler.TrackImpression)

	// Watch page & streaming
	app.Get("/watch/:token", streamHandler.ShowPlayer)
	app.Get("/stream/:token", middleware.RateLimitStream(), streamHandler.StreamVideo)

	// Admin routes
	admin := api.Group("/admin", middleware.AuthRequired(), middleware.AdminOnly())

	// Admin video management
	adminVideos := admin.Group("/videos")
	adminVideos.Get("/", adminVideoHandler.GetAllVideos)
	adminVideos.Post("/upload", adminVideoHandler.UploadVideo)
	adminVideos.Put("/:id", adminVideoHandler.UpdateVideo)
	adminVideos.Delete("/:id", adminVideoHandler.DeleteVideo)
	adminVideos.Post("/:id/refresh", adminVideoHandler.RefreshVideoLink)

	// Admin ads management
	adminAds := admin.Group("/ads")
	adminAds.Get("/", adminAdHandler.GetAllAds)
	adminAds.Get("/:id", adminAdHandler.GetAdByID)
	adminAds.Post("/", adminAdHandler.CreateAd)
	adminAds.Put("/:id", adminAdHandler.UpdateAd)
	adminAds.Delete("/:id", adminAdHandler.DeleteAd)
	adminAds.Post("/:id/toggle", adminAdHandler.ToggleActive)

	// Admin analytics
	adminAnalytics := admin.Group("/analytics")
	adminAnalytics.Get("/overview", analyticsHandler.GetOverview)
	adminAnalytics.Get("/daily", analyticsHandler.GetDailyStats)
	adminAnalytics.Get("/range", analyticsHandler.GetStatsByDateRange)
	adminAnalytics.Get("/monthly", analyticsHandler.GetMonthlyStats)
	adminAnalytics.Get("/top-videos", analyticsHandler.GetTopVideos)

	// Admin pCloud accounts management
	adminPCloud := admin.Group("/pcloud/accounts")
	adminPCloud.Get("/", pcloudHandler.GetAllAccounts)
	adminPCloud.Get("/:id", pcloudHandler.GetAccountByID)
	adminPCloud.Post("/", pcloudHandler.CreateAccount)
	adminPCloud.Put("/:id", pcloudHandler.UpdateAccount)
	adminPCloud.Delete("/:id", pcloudHandler.DeleteAccount)
	adminPCloud.Post("/:id/toggle", pcloudHandler.ToggleActive)

	// Initialize cron jobs
	c := cronpkg.New()
	
	// Refresh expired pCloud links (every 6 hours)
	refreshLinksJob := cron.NewRefreshLinksJob(pcloudService)
	c.AddFunc(config.GlobalConfig.Cron.RefreshLinks, refreshLinksJob.Run)

	// Aggregate daily stats (every midnight)
	aggregateStatsJob := cron.NewAggregateStatsJob(analyticsService)
	c.AddFunc(config.GlobalConfig.Cron.AggregateStats, aggregateStatsJob.Run)

	c.Start()
	log.Println("✅ Cron jobs started")

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("🛑 Shutting down server...")
		c.Stop()
		app.Shutdown()
	}()

	// Start server
	port := config.GlobalConfig.App.Port
	log.Printf("🚀 Server starting on port %s\n", port)
	log.Printf("🌐 Environment: %s\n", config.GlobalConfig.App.Env)
	
	if err := app.Listen(fmt.Sprintf(":%s", port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}