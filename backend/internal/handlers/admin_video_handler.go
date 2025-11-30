package handlers

import (
	"bobastream/internal/models"
	"bobastream/internal/services"
	"bobastream/internal/utils"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AdminVideoHandler struct {
	videoService    *services.VideoService
	pcloudService   *services.PCloudService
	categoryService *services.CategoryService
}

func NewAdminVideoHandler(
	videoService *services.VideoService,
	pcloudService *services.PCloudService,
	categoryService *services.CategoryService,
) *AdminVideoHandler {
	return &AdminVideoHandler{
		videoService:    videoService,
		pcloudService:   pcloudService,
		categoryService: categoryService,
	}
}

// UploadVideo uploads video to pCloud with auto-rotate (admin)
func (h *AdminVideoHandler) UploadVideo(c *fiber.Ctx) error {
	// Parse multipart form
	file, err := c.FormFile("video")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Video file is required")
	}

	// ✅ SANITIZE ALL USER INPUTS
	title := utils.SanitizeString(c.FormValue("title"))
	description := utils.SanitizeString(c.FormValue("description"))
	thumbnailURL := utils.SanitizeURL(c.FormValue("thumbnail_url"))
	categoryIDStr := c.FormValue("category_id")
	tagsStr := c.FormValue("tags") // comma-separated
	durationStr := c.FormValue("duration_seconds")

	// Validate required fields
	if title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Title is required")
	}

	// ✅ TRUNCATE TITLE (max 500 chars as per DB schema)
	title = utils.TruncateString(title, 500)

	// ✅ TRUNCATE DESCRIPTION (max 10000 chars reasonable limit)
	description = utils.TruncateString(description, 10000)

	// Parse duration (admin manual input)
	durationSeconds := 0
	if durationStr != "" {
		durationSeconds, _ = strconv.Atoi(durationStr)
		if durationSeconds < 0 {
			durationSeconds = 0
		}
	}

	// Open file
	fileHandle, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to open file")
	}
	defer fileHandle.Close()

	// Upload to pCloud (auto-rotate based on storage)
	credential, fileID, _, err := h.pcloudService.UploadFile(fileHandle, file.Filename, file.Size)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to upload to pCloud: "+err.Error())
	}

	// Get streaming link from pCloud
	sourceURL, expiresAt, err := h.pcloudService.GetFileLink(fileID, credential.APIToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get pCloud link")
	}

	// Parse category ID
	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		id, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &id
		}
	}

	// ✅ SANITIZE TAGS
	var tags []string
	if tagsStr != "" {
		rawTags := strings.Split(tagsStr, ",")
		tags = utils.SanitizeStrings(rawTags)
		
		// ✅ LIMIT NUMBER OF TAGS (max 20)
		if len(tags) > 20 {
			tags = tags[:20]
		}
		
		// ✅ LIMIT TAG LENGTH (max 50 chars each)
		for i := range tags {
			tags[i] = utils.TruncateString(tags[i], 50)
		}
	}

	// Create video record
	video := &models.Video{
		Title:              title,
		Description:        description,
		ThumbnailURL:       thumbnailURL,
		SourceURL:          sourceURL,
		SourceURLExpiresAt: &expiresAt,
		DurationSeconds:    durationSeconds,
		FileSizeMB:         float64(file.Size) / (1024 * 1024),
		PCloudFileID:       fmt.Sprintf("%d", fileID),
		PCloudCredentialID: credential.ID,
		CategoryID:         categoryID,
		Tags:               tags,
		IsPublished:        true,
	}

	// Save video to database
	if err := h.videoService.CreateVideo(video); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to save video")
	}

	// Create wrapper link
	wrapperLink := &models.WrapperLink{
		VideoID: video.ID,
	}
	if err := h.videoService.CreateWrapperLink(wrapperLink); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create wrapper link")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"video":        video,
		"wrapper_link": wrapperLink,
		"pcloud_account": fiber.Map{
			"account_name":  credential.AccountName,
			"storage_used":  credential.StorageUsedGB,
			"storage_limit": credential.StorageLimitGB,
		},
	}, "Video uploaded successfully")
}

// GetAllVideos gets all videos (admin)
func (h *AdminVideoHandler) GetAllVideos(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	// ✅ VALIDATE PAGINATION PARAMS
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	videos, total, err := h.videoService.GetAllVideos(page, limit)
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

// UpdateVideo updates video (admin)
func (h *AdminVideoHandler) UpdateVideo(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	video, err := h.videoService.GetVideoByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Video not found")
	}

	var req struct {
		Title        string   `json:"title"`
		Description  string   `json:"description"`
		ThumbnailURL string   `json:"thumbnail_url"`
		CategoryID   *string  `json:"category_id"`
		Tags         []string `json:"tags"`
		IsPublished  *bool    `json:"is_published"`
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
		video.Title = utils.TruncateString(sanitizedTitle, 500)
	}

	if req.Description != "" {
		video.Description = utils.TruncateString(utils.SanitizeString(req.Description), 10000)
	}

	if req.ThumbnailURL != "" {
		sanitizedURL := utils.SanitizeURL(req.ThumbnailURL)
		if sanitizedURL == "" {
			return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid thumbnail URL")
		}
		video.ThumbnailURL = sanitizedURL
	}

	if req.CategoryID != nil {
		if *req.CategoryID == "" {
			video.CategoryID = nil
		} else {
			catID, err := uuid.Parse(*req.CategoryID)
			if err == nil {
				video.CategoryID = &catID
			}
		}
	}

	if req.Tags != nil {
		// ✅ SANITIZE TAGS
		sanitizedTags := utils.SanitizeStrings(req.Tags)
		
		// ✅ LIMIT NUMBER OF TAGS (max 20)
		if len(sanitizedTags) > 20 {
			sanitizedTags = sanitizedTags[:20]
		}
		
		// ✅ LIMIT TAG LENGTH (max 50 chars each)
		for i := range sanitizedTags {
			sanitizedTags[i] = utils.TruncateString(sanitizedTags[i], 50)
		}
		
		video.Tags = sanitizedTags
	}

	if req.IsPublished != nil {
		video.IsPublished = *req.IsPublished
		if *req.IsPublished && video.PublishedAt == nil {
			now := time.Now()
			video.PublishedAt = &now
		}
	}

	if err := h.videoService.UpdateVideo(video); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update video")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"video": video,
	}, "Video updated successfully")
}

// DeleteVideo deletes video (admin)
func (h *AdminVideoHandler) DeleteVideo(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	if err := h.videoService.DeleteVideo(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete video")
	}

	return utils.SuccessResponse(c, nil, "Video deleted successfully")
}

// RefreshVideoLink manually refreshes video pCloud link (admin)
func (h *AdminVideoHandler) RefreshVideoLink(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid video ID")
	}

	video, err := h.videoService.GetVideoByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Video not found")
	}

	if video.PCloudCredential == nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "No pCloud credential associated")
	}

	// Get new link
	fileID := int64(0)
	fmt.Sscanf(video.PCloudFileID, "%d", &fileID)

	newURL, expiresAt, err := h.pcloudService.GetFileLink(fileID, video.PCloudCredential.APIToken)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to refresh link")
	}

	// Update video
	if err := h.videoService.UpdateSourceURL(id, newURL, expiresAt); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update video")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"source_url":            newURL,
		"source_url_expires_at": expiresAt,
	}, "Link refreshed successfully")
}