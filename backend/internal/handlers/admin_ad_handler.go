package handlers

import (
	"bobastream/internal/models"
	"bobastream/internal/services"
	"bobastream/internal/utils"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminAdHandler struct {
	adService *services.AdService
}

func NewAdminAdHandler(adService *services.AdService) *AdminAdHandler {
	return &AdminAdHandler{adService: adService}
}

// GetAllAds gets all ads (admin)
func (h *AdminAdHandler) GetAllAds(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// ✅ VALIDATE PAGINATION PARAMS
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	ads, total, err := h.adService.GetAllAds(page, limit)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get ads")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ads":   ads,
		"total": total,
		"page":  page,
		"limit": limit,
	}, "")
}

// GetAdByID gets ad by ID (admin)
func (h *AdminAdHandler) GetAdByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad ID")
	}

	ad, err := h.adService.GetAdByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Ad not found")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ad": ad,
	}, "")
}

// CreateAd creates a new ad (admin)
func (h *AdminAdHandler) CreateAd(c *fiber.Ctx) error {
	var req struct {
		Title            string `json:"title" validate:"required"`
		AdType           string `json:"ad_type" validate:"required"`
		ContentURL       string `json:"content_url"`
		RedirectURL      string `json:"redirect_url"`
		DurationSeconds  int    `json:"duration_seconds"`
		DisplayFrequency int    `json:"display_frequency"`
		Priority         int    `json:"priority"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// ✅ SANITIZE INPUTS
	title := utils.SanitizeString(req.Title)
	if title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Title is required")
	}
	title = utils.TruncateString(title, 255)

	// ✅ SANITIZE URLs
	contentURL := utils.SanitizeURL(req.ContentURL)
	redirectURL := utils.SanitizeURL(req.RedirectURL)

	// ✅ VALIDATE AD TYPE
	adType := models.AdTypePreroll
	switch req.AdType {
	case "banner":
		adType = models.AdTypeBanner
	case "popup":
		adType = models.AdTypePopup
	case "preroll":
		adType = models.AdTypePreroll
	default:
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad type. Must be: preroll, banner, or popup")
	}

	// ✅ VALIDATE NUMERIC FIELDS
	if req.DurationSeconds < 0 {
		req.DurationSeconds = 0
	}
	if req.DurationSeconds > 300 { // Max 5 minutes
		req.DurationSeconds = 300
	}

	if req.DisplayFrequency < 1 {
		req.DisplayFrequency = 1
	}
	if req.DisplayFrequency > 100 {
		req.DisplayFrequency = 100
	}

	if req.Priority < 0 {
		req.Priority = 0
	}
	if req.Priority > 999 {
		req.Priority = 999
	}

	ad := &models.Ad{
		Title:            title,
		AdType:           adType,
		ContentURL:       contentURL,
		RedirectURL:      redirectURL,
		DurationSeconds:  req.DurationSeconds,
		DisplayFrequency: req.DisplayFrequency,
		Priority:         req.Priority,
		IsActive:         true,
	}

	if err := h.adService.CreateAd(ad); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create ad")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ad": ad,
	}, "Ad created successfully")
}

// UpdateAd updates an ad (admin)
func (h *AdminAdHandler) UpdateAd(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad ID")
	}

	ad, err := h.adService.GetAdByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Ad not found")
	}

	var req struct {
		Title            string `json:"title"`
		AdType           string `json:"ad_type"`
		ContentURL       string `json:"content_url"`
		RedirectURL      string `json:"redirect_url"`
		DurationSeconds  int    `json:"duration_seconds"`
		DisplayFrequency int    `json:"display_frequency"`
		Priority         int    `json:"priority"`
		IsActive         *bool  `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// ✅ SANITIZE AND UPDATE FIELDS
	if req.Title != "" {
		sanitizedTitle := utils.SanitizeString(req.Title)
		if sanitizedTitle == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Title cannot be empty")
		}
		ad.Title = utils.TruncateString(sanitizedTitle, 255)
	}

	if req.ContentURL != "" {
		sanitizedURL := utils.SanitizeURL(req.ContentURL)
		if sanitizedURL == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid content URL")
		}
		ad.ContentURL = sanitizedURL
	}

	if req.RedirectURL != "" {
		sanitizedURL := utils.SanitizeURL(req.RedirectURL)
		if sanitizedURL == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid redirect URL")
		}
		ad.RedirectURL = sanitizedURL
	}

	// ✅ VALIDATE NUMERIC FIELDS
	if req.DurationSeconds > 0 {
		if req.DurationSeconds > 300 {
			req.DurationSeconds = 300
		}
		ad.DurationSeconds = req.DurationSeconds
	}

	if req.DisplayFrequency > 0 {
		if req.DisplayFrequency > 100 {
			req.DisplayFrequency = 100
		}
		ad.DisplayFrequency = req.DisplayFrequency
	}

	if req.Priority >= 0 {
		if req.Priority > 999 {
			req.Priority = 999
		}
		ad.Priority = req.Priority
	}

	if req.IsActive != nil {
		ad.IsActive = *req.IsActive
	}

	if err := h.adService.UpdateAd(ad); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update ad")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"ad": ad,
	}, "Ad updated successfully")
}

// DeleteAd deletes an ad (admin)
func (h *AdminAdHandler) DeleteAd(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad ID")
	}

	if err := h.adService.DeleteAd(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete ad")
	}

	return utils.SuccessResponse(c, nil, "Ad deleted successfully")
}

// ToggleActive toggles ad active status (admin)
func (h *AdminAdHandler) ToggleActive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid ad ID")
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.adService.ToggleActive(id, req.IsActive); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to toggle ad status")
	}

	return utils.SuccessResponse(c, nil, "Ad status updated")
}