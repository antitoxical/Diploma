package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetAllDeliveryMethods(c *fiber.Ctx) error {
	methods, err := h.deliveryService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_delivery_methods", err.Error()))
	}
	if methods == nil {
		methods = []models.DeliveryMethod{}
	}
	return c.JSON(web.OkResp("success_delivery_methods_retrieved", methods))
}

func (h *Handler) CreateDeliveryMethod(c *fiber.Ctx) error {
	var input models.DeliveryMethod
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	m, err := h.deliveryService.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_delivery_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_delivery_method_created", m))
}

func (h *Handler) UpdateDeliveryMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid ID"))
	}

	var input models.DeliveryMethod
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.deliveryService.Update(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_delivery_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_delivery_method_updated", nil))
}

func (h *Handler) DeleteDeliveryMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid ID"))
	}

	if err := h.deliveryService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_delivery_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_delivery_method_deleted", nil))
}
