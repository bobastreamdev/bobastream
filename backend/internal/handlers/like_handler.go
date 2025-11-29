package handlers

import (
	"bobastream/internal/services"
	"bobastream/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LikeHandler struct {
	videoService *services.VideoService
}

func NewLikeHandler(videoService *services.VideoService) *LikeHandler {
	return &LikeHandler{videoService: videoService}
}

// LikeVideo likes a video
func (h *LikeHandler) LikeVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
	}

	if err := h.videoService.LikeVideo(videoID, userID); err != nil {
		if err.Error() == "already liked" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Already liked")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to like video")
	}

	return utils.SuccessResponse(c, nil, "Video liked")
}

// UnlikeVideo unlikes a video
func (h *LikeHandler) UnlikeVideo(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
	}

	if err := h.videoService.UnlikeVideo(videoID, userID); err != nil {
		if err.Error() == "not liked" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Not liked")
		}
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to unlike video")
	}

	return utils.SuccessResponse(c, nil, "Video unliked")
}

// CheckLiked checks if user has liked video
func (h *LikeHandler) CheckLiked(c *fiber.Ctx) error {
	videoID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.SuccessResponse(c, fiber.Map{
			"is_liked": false,
		}, "")
	}

	isLiked, err := h.videoService.IsVideoLiked(videoID, userID)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to check like status")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"is_liked": isLiked,
	}, "")
}

// GetUserLikedVideos gets user's liked videos
func (h *LikeHandler) GetUserLikedVideos(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uuid.UUID)
	if !ok {
		return utils.ErrorResponse(c, fiber.StatusUnauthorized, "Authentication required")
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	videos, total, err := h.videoService.GetUserLikedVideos(userID, page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get liked videos")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"videos": videos,
		"total":  total,
		"page":   page,
		"limit":  limit,
	}, "")
}