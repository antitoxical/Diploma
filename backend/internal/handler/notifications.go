package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CreateNotification godoc
// @Summary Create notification
// @Description Creates a new notification for a user (admin)
// @Tags notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body models.CreateNotification true "Notification data"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /admin/notifications [post]
func (h *Handler) CreateNotification(c *fiber.Ctx) error {
	var input models.CreateNotification
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	n, err := h.notificationService.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_notification", err.Error()))
	}

	return c.JSON(web.OkResp("success_notification_created", n))
}

// GetUserNotifications godoc
// @Summary Get user notifications
// @Description Returns all notifications for a user
// @Tags notifications
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /notifications/{user_id} [get]
func (h *Handler) GetUserNotifications(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	list, err := h.notificationService.GetByUser(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_notifications", err.Error()))
	}

	return c.JSON(web.OkResp("success_notifications_retrieved", list))
}

// MarkNotificationRead godoc
// @Summary Mark notification as read
// @Description Marks a single notification as read
// @Tags notifications
// @Produce json
// @Param id path string true "Notification UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /notifications/{id}/read [patch]
func (h *Handler) MarkNotificationRead(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid notification ID"))
	}

	if err := h.notificationService.MarkRead(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_mark_read", err.Error()))
	}

	return c.JSON(web.OkResp("success_notification_marked_read", nil))
}

// MarkAllNotificationsRead godoc
// @Summary Mark all notifications as read
// @Description Marks all notifications for a user as read
// @Tags notifications
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /notifications/{user_id}/read-all [patch]
func (h *Handler) MarkAllNotificationsRead(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	if err := h.notificationService.MarkAllRead(c.Context(), telegramID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_mark_all_read", err.Error()))
	}

	return c.JSON(web.OkResp("success_all_notifications_marked_read", nil))
}

// DeleteNotification godoc
// @Summary Delete notification
// @Description Deletes a notification (admin)
// @Tags notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "Notification UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /admin/notifications/{id} [delete]
func (h *Handler) DeleteNotification(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid notification ID"))
	}

	if err := h.notificationService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_notification", err.Error()))
	}

	return c.JSON(web.OkResp("success_notification_deleted", nil))
}
