package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetAllPaymentMethods(c *fiber.Ctx) error {
	methods, err := h.paymentService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_payment_methods", err.Error()))
	}
	if methods == nil {
		methods = []models.PaymentMethod{}
	}
	return c.JSON(web.OkResp("success_payment_methods_retrieved", methods))
}

func (h *Handler) CreatePaymentMethod(c *fiber.Ctx) error {
	var input models.PaymentMethod
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	m, err := h.paymentService.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_payment_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_payment_method_created", m))
}

func (h *Handler) UpdatePaymentMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid ID"))
	}

	var input models.PaymentMethod
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.paymentService.Update(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_payment_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_payment_method_updated", nil))
}

func (h *Handler) DeletePaymentMethod(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid ID"))
	}

	if err := h.paymentService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_payment_method", err.Error()))
	}
	return c.JSON(web.OkResp("success_payment_method_deleted", nil))
}
