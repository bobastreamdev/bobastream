package handlers

import (
	"bobastream/internal/services"
	"bobastream/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type VideoHandler struct {
	videoService    *services.VideoService
	categoryService *services.CategoryService
}

func NewVideoHandler(
	videoService *services.VideoService,
	categoryService *services.CategoryService,
) *VideoHandler {
	return &VideoHandler{
		videoService:    videoService,
		categoryService: categoryService,
	}
}

// GetFeed gets video feed
func (h *VideoHandler) GetFeed(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// ✅ VALIDATE PAGINATION
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	videos, total, err := h.videoService.GetFeedVideos(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}, "")
}

// GetVideoByID gets video by ID
func (h *VideoHandler) GetVideoByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	video, err := h.videoService.GetVideoByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Video not found")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"video": video,
	}, "")
}

// GetRelatedVideos gets related videos
func (h *VideoHandler) GetRelatedVideos(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	// ✅ VALIDATE LIMIT
	if limit < 1 || limit > 50 {
		limit = 10
	}

	videos, err := h.videoService.GetRelatedVideos(id, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get related videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
	}, "")
}

// ✅ FIXED: SearchVideos with keyword validation and sanitization
func (h *VideoHandler) SearchVideos(c *fiber.Ctx) error {
	keyword := c.Query("q", "")
	categoryIDStr := c.Query("category_id", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// ✅ VALIDATE AND SANITIZE KEYWORD
	if len(keyword) > 100 {
		keyword = keyword[:100]
	}
	keyword = utils.SanitizeString(keyword)

	// ✅ VALIDATE PAGINATION
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &id
		}
	}

	videos, total, err := h.videoService.SearchVideos(keyword, categoryID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to search videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}, "")
}

// GetVideosByCategory gets videos by category
func (h *VideoHandler) GetVideosByCategory(c *fiber.Ctx) error {
	categoryID, err := uuid.Parse(c.Params("categoryId"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid category ID")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// ✅ VALIDATE PAGINATION
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	videos, total, err := h.videoService.GetVideosByCategory(categoryID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}, "")
}

// GetCategories gets all active categories
func (h *VideoHandler) GetCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetActiveCategories()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get categories")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"categories": categories,
	}, "")
}

// TrackView tracks video view with watch duration
func (h *VideoHandler) TrackView(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	var req struct {
		SessionID     string `json:"session_id" validate:"required"`
		WatchDuration int    `json:"watch_duration" validate:"required"`
		VideoDuration int    `json:"video_duration" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// ✅ VALIDATE DURATIONS
	if req.WatchDuration < 0 {
		req.WatchDuration = 0
	}
	if req.VideoDuration < 0 {
		req.VideoDuration = 0
	}

	// Get user ID if authenticated
	var userID *uuid.UUID
	if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &uid
	}

	viewerIP := c.IP()
	userAgent := c.Get("User-Agent")

	err = h.videoService.TrackVideoView(
		id,
		userID,
		req.SessionID,
		viewerIP,
		userAgent,
		req.WatchDuration,
		req.VideoDuration,
	)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to track view")
	}

	return utils.SuccessResponse(c, nil, "View tracked")
}