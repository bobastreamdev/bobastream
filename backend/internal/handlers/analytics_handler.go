package handlers

import (
	"bobastream/internal/services"
	"bobastream/internal/utils"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	analyticsService *services.AnalyticsService
	videoService     *services.VideoService
}

func NewAnalyticsHandler(
	analyticsService *services.AnalyticsService,
	videoService *services.VideoService,
) *AnalyticsHandler {
	return &AnalyticsHandler{
		analyticsService: analyticsService,
		videoService:     videoService,
	}
}

// GetOverview gets dashboard overview stats (admin)
func (h *AnalyticsHandler) GetOverview(c *fiber.Ctx) error {
	stats, err := h.analyticsService.GetOverviewStats()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get overview stats")
	}

	return utils.SuccessResponse(c, stats, "")
}

// GetDailyStats gets daily stats (admin)
func (h *AnalyticsHandler) GetDailyStats(c *fiber.Ctx) error {
	dateStr := c.Query("date", time.Now().Format("2006-01-02"))
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid date format")
	}

	stats, err := h.analyticsService.GetDailyStats(date)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Stats not found for this date")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"stats": stats,
	}, "")
}

// GetStatsByDateRange gets stats within date range (admin)
func (h *AnalyticsHandler) GetStatsByDateRange(c *fiber.Ctx) error {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "start_date and end_date are required")
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid start_date format")
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid end_date format")
	}

	stats, err := h.analyticsService.GetStatsByDateRange(startDate, endDate)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get stats")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"stats":      stats,
		"start_date": startDate,
		"end_date":   endDate,
	}, "")
}

// GetMonthlyStats gets aggregated monthly stats (admin)
func (h *AnalyticsHandler) GetMonthlyStats(c *fiber.Ctx) error {
	year, _ := strconv.Atoi(c.Query("year", strconv.Itoa(time.Now().Year())))
	month, _ := strconv.Atoi(c.Query("month", strconv.Itoa(int(time.Now().Month()))))

	if year < 2020 || year > 2100 || month < 1 || month > 12 {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid year or month")
	}

	stats, err := h.analyticsService.GetMonthlyStats(year, month)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get monthly stats")
	}

	return utils.SuccessResponse(c, stats, "")
}

// GetTopVideos gets top performing videos (admin)
func (h *AnalyticsHandler) GetTopVideos(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	days, _ := strconv.Atoi(c.Query("days", "7"))

	videos, err := h.videoService.GetTopVideos(limit, days)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get top videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
		"limit":  limit,
		"days":   days,
	}, "")
}