package handlers

import (
	"bobastream/internal/models"
	"bobastream/internal/services"
	"bobastream/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type PCloudHandler struct {
	pcloudService *services.PCloudService
}

func NewPCloudHandler(pcloudService *services.PCloudService) *PCloudHandler {
	return &PCloudHandler{pcloudService: pcloudService}
}

// GetAllAccounts gets all pCloud accounts (admin)
func (h *PCloudHandler) GetAllAccounts(c *fiber.Ctx) error {
	accounts, err := h.pcloudService.GetAllCredentials()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to get accounts")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"accounts": accounts,
	}, "")
}

// GetAccountByID gets pCloud account by ID (admin)
func (h *PCloudHandler) GetAccountByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid account ID")
	}

	account, err := h.pcloudService.GetCredentialByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Account not found")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"account": account,
	}, "")
}

// CreateAccount creates new pCloud account (admin)
func (h *PCloudHandler) CreateAccount(c *fiber.Ctx) error {
	var req struct {
		AccountName    string  `json:"account_name" validate:"required"`
		APIToken       string  `json:"api_token" validate:"required"`
		StorageLimitGB float64 `json:"storage_limit_gb" validate:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	credential := &models.PCloudCredential{
		AccountName:    req.AccountName,
		APIToken:       req.APIToken,
		StorageLimitGB: req.StorageLimitGB,
		StorageUsedGB:  0,
		IsActive:       true,
	}

	if err := h.pcloudService.CreateCredential(credential); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to create account")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"account": credential,
	}, "Account created successfully")
}

// UpdateAccount updates pCloud account (admin)
func (h *PCloudHandler) UpdateAccount(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid account ID")
	}

	account, err := h.pcloudService.GetCredentialByID(id)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusNotFound, "Account not found")
	}

	var req struct {
		AccountName    string   `json:"account_name"`
		APIToken       string   `json:"api_token"`
		StorageLimitGB *float64 `json:"storage_limit_gb"`
		StorageUsedGB  *float64 `json:"storage_used_gb"`
		IsActive       *bool    `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Update fields
	if req.AccountName != "" {
		account.AccountName = req.AccountName
	}
	if req.APIToken != "" {
		account.APIToken = req.APIToken
	}
	if req.StorageLimitGB != nil {
		account.StorageLimitGB = *req.StorageLimitGB
	}
	if req.StorageUsedGB != nil {
		account.StorageUsedGB = *req.StorageUsedGB
	}
	if req.IsActive != nil {
		account.IsActive = *req.IsActive
	}

	if err := h.pcloudService.UpdateCredential(account); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to update account")
	}

	return utils.SuccessResponse(c, fiber.Map{
		"account": account,
	}, "Account updated successfully")
}

// DeleteAccount deletes pCloud account (admin)
func (h *PCloudHandler) DeleteAccount(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid account ID")
	}

	if err := h.pcloudService.DeleteCredential(id); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to delete account")
	}

	return utils.SuccessResponse(c, nil, "Account deleted successfully")
}

// ToggleActive toggles account active status (admin)
func (h *PCloudHandler) ToggleActive(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid account ID")
	}

	var req struct {
		IsActive bool `json:"is_active"`
	}

	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	if err := h.pcloudService.ToggleActive(id, req.IsActive); err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "Failed to toggle account status")
	}

	return utils.SuccessResponse(c, nil, "Account status updated")
}