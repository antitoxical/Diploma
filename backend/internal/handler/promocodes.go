package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetAllPromocodes(c *fiber.Ctx) error {
	promos, err := h.promocodeService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_promocodes", err.Error()))
	}
	if promos == nil {
		promos = []models.Promocode{}
	}
	return c.JSON(web.OkResp("success_promocodes_retrieved", promos))
}

func (h *Handler) CreatePromocode(c *fiber.Ctx) error {
	var input models.Promocode
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	promo, err := h.promocodeService.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_promocode", err.Error()))
	}
	return c.JSON(web.OkResp("success_promocode_created", promo))
}

func (h *Handler) DeletePromocode(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid ID"))
	}

	if err := h.promocodeService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_promocode", err.Error()))
	}
	return c.JSON(web.OkResp("success_promocode_deleted", nil))
}

func (h *Handler) ValidatePromocode(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_missing_code", "code query parameter required"))
	}

	promo, err := h.promocodeService.ValidateCode(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.ErrorResp("error_invalid_promo", "Promocode not found or expired"))
	}
	return c.JSON(web.OkResp("success_promocode_valid", promo))
}
