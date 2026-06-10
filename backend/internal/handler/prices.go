package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreatePrice(c *fiber.Ctx) error {
	var input models.ProductPriceTier
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}
	if input.MinQuantity < 1 || input.PriceKopeks < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_data", "Invalid min_quantity or price_kopeks"))
	}

	tier, err := h.priceService.CreateTier(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_price", err.Error()))
	}

	return c.JSON(web.OkResp("success_price_created", tier))
}

func (h *Handler) GetPriceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid price tier ID"))
	}

	tier, err := h.priceService.GetTierByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_price", err.Error()))
	}

	return c.JSON(web.OkResp("success_price_retrieved", tier))
}

func (h *Handler) GetPricesByProductID(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	tiers, err := h.priceService.GetTiersByProductID(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_prices_by_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_prices_retrieved", tiers))
}

func (h *Handler) UpdatePrice(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid price tier ID"))
	}
	var input models.UpdatePriceTierInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.priceService.UpdateTier(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_price", err.Error()))
	}

	return c.JSON(web.OkResp("success_price_updated", nil))
}

func (h *Handler) DeletePrice(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid price tier ID"))
	}

	if err := h.priceService.DeleteTier(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_price", err.Error()))
	}

	return c.JSON(web.OkResp("success_price_deleted", nil))
}

func (h *Handler) DeletePricesByProductID(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	if err := h.priceService.DeleteTiersByProductID(c.Context(), productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_prices_by_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_prices_deleted", nil))
}

func (h *Handler) UpdatePriceCount(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid price tier ID"))
	}

	var input models.UpdateTierMinQuantity
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.priceService.UpdateTierMinQuantity(c.Context(), id, input.NewMinQuantity); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_price_count", err.Error()))
	}

	return c.JSON(web.OkResp("success_price_count_updated", nil))
}
