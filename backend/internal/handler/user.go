package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var input models.CreateUser
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	user, err := h.userService.CreateUser(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_user", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_created", user))
}

func (h *Handler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	telegramID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	user, err := h.userService.GetUserByTelegramID(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.ErrorResp("error_get_user", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_retrieved", user))
}

func (h *Handler) UpdateUserProfile(c *fiber.Ctx) error {
	idStr := c.Params("id")
	telegramID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	var input models.UpdateUserProfile
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.userService.UpdateProfile(c.Context(), telegramID, input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_update_user", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_updated", nil))
}

func (h *Handler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	telegramID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	if err := h.userService.DeleteUser(c.Context(), telegramID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_user", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_deleted", nil))
}

func (h *Handler) GetAllUsers(c *fiber.Ctx) error {
	users, err := h.userService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_all_users", err.Error()))
	}
	return c.JSON(web.OkResp("success_all_users_retrieved", users))
}

// ── User Addresses ──────────────────────────────────────────────────

func (h *Handler) CreateUserAddress(c *fiber.Ctx) error {
	var input models.CreateUserAddress
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	addr, err := h.userService.CreateAddress(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_create_address", err.Error()))
	}
	return c.JSON(web.OkResp("success_address_created", addr))
}

func (h *Handler) GetUserAddresses(c *fiber.Ctx) error {
	telegramID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	addrs, err := h.userService.GetAddresses(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_addresses", err.Error()))
	}
	return c.JSON(web.OkResp("success_addresses_retrieved", addrs))
}

func (h *Handler) DeleteUserAddress(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid address ID"))
	}

	if err := h.userService.DeleteAddress(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_address", err.Error()))
	}
	return c.JSON(web.OkResp("success_address_deleted", nil))
}
