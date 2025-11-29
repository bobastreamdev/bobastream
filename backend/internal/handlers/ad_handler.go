package handlers

import (
	"bobastream/internal/models"
	"bobastream/internal/services"
	"bobastream/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdHandler struct {
	adService *services.AdService
}

func NewAdHandler(adService *services.AdService) *AdHandler {
	return &AdHandler{adService: adService}
}

// GetPrerollAd gets active preroll ad
func (h *AdHandler) GetPrerollAd(c *fiber.Ctx) error {
	ad, err := h.adService.GetActiveAdByType(models.AdTypePreroll)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No active ad available")
	}

	return utils.SuccessResponse(c, ad, "")
}

// GetBannerAd gets active banner ad
func (h *AdHandler) GetBannerAd(c *fiber.Ctx) error {
	ad, err := h.adService.GetActiveAdByType(models.AdTypeBanner)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No active ad available")
	}

	return utils.SuccessResponse(c, ad, "")
}

// GetPopupAd gets active popup ad
func (h *AdHandler) GetPopupAd(c *fiber.Ctx) error {
	ad, err := h.adService.GetActiveAdByType(models.AdTypePopup)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "No active ad available")
	}

	return utils.SuccessResponse(c, ad, "")
}

// TrackImpression tracks ad impression (view/click/skip)
func (h *AdHandler) TrackImpression(c *fiber.Ctx) error {
	adID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad ID")
	}

	var req struct {
		ImpressionType  string     `json:"impression_type" validate:"required"`
		WatchedDuration int        `json:"watched_duration"`
		SessionID       string     `json:"session_id"`
		VideoID         *uuid.UUID `json:"video_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get user ID if authenticated
	var userID *uuid.UUID
	if uid, ok := c.Locals("user_id").(uuid.UUID); ok {
		userID = &uid
	}

	viewerIP := c.IP()

	// Parse impression type
	impressionType := models.ImpressionView
	switch req.ImpressionType {
	case "click":
		impressionType = models.ImpressionClick
	case "skip":
		impressionType = models.ImpressionSkip
	}

	err = h.adService.TrackAdImpression(
		adID,
		req.VideoID,
		userID,
		viewerIP,
		impressionType,
		req.WatchedDuration,
		req.SessionID,
	)

	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to track impression")
	}

	return utils.SuccessResponse(c, nil, "Impression tracked")
}