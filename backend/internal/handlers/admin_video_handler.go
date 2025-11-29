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

	title := c.FormValue("title")
	description := c.FormValue("description")
	thumbnailURL := c.FormValue("thumbnail_url")
	categoryIDStr := c.FormValue("category_id")
	tagsStr := c.FormValue("tags") // comma-separated

	if title == "" {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Title is required")
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

	// Parse tags
	var tags []string
	if tagsStr != "" {
		tags = strings.Split(tagsStr, ",")
		// Trim spaces
		for i := range tags {
			tags[i] = strings.TrimSpace(tags[i])
		}
	}

	// Create video record
	video := &models.Video{
		Title:              title,
		Description:        description,
		ThumbnailURL:       thumbnailURL,
		SourceURL:          sourceURL,
		SourceURLExpiresAt: &expiresAt,
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
			"account_name": credential.AccountName,
			"storage_used": credential.StorageUsedGB,
			"storage_limit": credential.StorageLimitGB,
		},
	}, "Video uploaded successfully")
}

// GetAllVideos gets all videos (admin)
func (h *AdminVideoHandler) GetAllVideos(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

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

	// Update fields
	if req.Title != "" {
		video.Title = req.Title
	}
	if req.Description != "" {
		video.Description = req.Description
	}
	if req.ThumbnailURL != "" {
		video.ThumbnailURL = req.ThumbnailURL
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
		video.Tags = req.Tags
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
		"source_url":           newURL,
		"source_url_expires_at": expiresAt,
	}, "Link refreshed successfully")
}