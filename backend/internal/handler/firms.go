package handler

import (
	"database/sql"
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateFirm godoc
// @Summary Create firm
// @Description Creates a new firm/brand (admin)
// @Tags firms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.Firm true "Firm data"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /admin/firms [post]
func (h *Handler) CreateFirm(c *fiber.Ctx) error {
	var input models.Firm
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	firm, err := h.firmsService.CreateFirm(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_firm", err.Error()))
	}

	return c.JSON(web.OkResp("success_firm_created", firm))
}

// GetFirmByID godoc
// @Summary Get firm by ID
// @Description Returns firm details by UUID
// @Tags firms
// @Produce json
// @Param id path string true "Firm UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 404 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /firms/{id} [get]
func (h *Handler) GetFirmByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid firm ID"))
	}

	firm, err := h.firmsService.GetFirmByID(c.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(web.ErrorResp("error_firm_not_found", "Firm not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_firm", err.Error()))
	}

	return c.JSON(web.OkResp("success_firm_retrieved", firm))
}

// GetAllFirms godoc
// @Summary Get all firms
// @Description Returns list of all firms/brands
// @Tags firms
// @Produce json
// @Success 200 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /firms [get]
func (h *Handler) GetAllFirms(c *fiber.Ctx) error {
	firms, err := h.firmsService.GetAllFirms(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_all_firms", err.Error()))
	}

	return c.JSON(web.OkResp("success_all_firms_retrieved", firms))
}

// UpdateFirm godoc
// @Summary Update firm
// @Description Updates firm details (admin)
// @Tags firms
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Firm UUID"
// @Param input body models.UpdateFirmInput true "Updated firm data"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /admin/firms/{id} [put]
func (h *Handler) UpdateFirm(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid firm ID"))
	}
	var input models.UpdateFirmInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.firmsService.UpdateFirm(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_firm", err.Error()))
	}

	return c.JSON(web.OkResp("success_firm_updated", nil))
}

// DeleteFirm godoc
// @Summary Delete firm
// @Description Deletes a firm (admin)
// @Tags firms
// @Produce json
// @Security BearerAuth
// @Param id path string true "Firm UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /admin/firms/{id} [delete]
func (h *Handler) DeleteFirm(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid firm ID"))
	}

	if err := h.firmsService.DeleteFirm(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_firm", err.Error()))
	}

	return c.JSON(web.OkResp("success_firm_deleted", nil))
}
