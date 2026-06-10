package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	var input models.CreateOrder
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	order, err := h.orderService.CreateOrder(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_create_order", err.Error()))
	}

	return c.JSON(web.OkResp("success_order_created", order))
}

func (h *Handler) GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_order_id", "Invalid order ID"))
	}

	order, err := h.orderService.GetOrderByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_order", err.Error()))
	}

	return c.JSON(web.OkResp("success_order_retrieved", order))
}

func (h *Handler) GetUserOrders(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	orders, err := h.orderService.GetUserOrders(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_user_orders", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_orders_retrieved", orders))
}

func (h *Handler) GetAllOrders(c *fiber.Ctx) error {
	orders, err := h.orderService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_all_orders", err.Error()))
	}
	if orders == nil {
		orders = []models.OrderWithProducts{}
	}
	return c.JSON(web.OkResp("success_all_orders_retrieved", orders))
}

func (h *Handler) UpdateOrderStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_order_id", "Invalid order ID"))
	}

	var input models.UpdateOrderStatus
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.orderService.UpdateStatus(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_update_order_status", err.Error()))
	}

	return c.JSON(web.OkResp("success_order_status_updated", nil))
}

func (h *Handler) GetOrderStatusHistory(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_order_id", "Invalid order ID"))
	}

	history, err := h.orderService.GetStatusHistory(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_status_history", err.Error()))
	}

	return c.JSON(web.OkResp("success_status_history_retrieved", history))
}
