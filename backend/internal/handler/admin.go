package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetAllAdmins(c *fiber.Ctx) error {
	admins, err := h.adminService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_admins", err.Error()))
	}
	if admins == nil {
		admins = []models.Admin{}
	}
	return c.JSON(web.OkResp("success_admins_retrieved", admins))
}

func (h *Handler) GetAdminByID(c *fiber.Ctx) error {
	id := c.Params("id")
	admin, err := h.adminService.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(web.ErrorResp("error_admin_not_found", "Admin not found"))
	}
	return c.JSON(web.OkResp("success_admin_retrieved", admin))
}

func (h *Handler) UpdateAdmin(c *fiber.Ctx) error {
	id := c.Params("id")
	var input models.AdminUpdateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}
	if input.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_missing_username", "Username is required"))
	}
	admin, err := h.adminService.UpdateAdmin(c.Context(), id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_admin", err.Error()))
	}
	return c.JSON(web.OkResp("success_admin_updated", admin))
}

func (h *Handler) AdminLogin(c *fiber.Ctx) error {
	var input models.AdminLoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}
	if input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_missing_credentials", "Username and password are required"))
	}

	resp, err := h.adminService.Login(c.Context(), input, h.jwtSecret)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(web.ErrorResp("error_invalid_credentials", "Invalid username or password"))
	}

	recordID := resp.Admin.ID
	h.writeAudit(c, resp.Admin.ID, "login", "admins", &recordID, map[string]any{
		"username": resp.Admin.Username,
	})

	return c.JSON(web.OkResp("success_admin_login", resp))
}

func (h *Handler) CreateAdmin(c *fiber.Ctx) error {
	var input models.AdminCreateInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}
	if input.Username == "" || input.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_missing_fields", "Username and password are required"))
	}

	admin, err := h.adminService.CreateAdmin(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_admin", err.Error()))
	}

	return c.JSON(web.OkResp("success_admin_created", admin))
}

func (h *Handler) GetAuditLogs(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "100"))

	logs, err := h.adminService.GetAuditLogs(c.Context(), limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_audit_logs", err.Error()))
	}
	if logs == nil {
		logs = []models.AuditLog{}
	}
	return c.JSON(web.OkResp("success_audit_logs_retrieved", logs))
}
