package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetAllTags(c *fiber.Ctx) error {
	tags, err := h.tagService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_tags", err.Error()))
	}
	return c.JSON(web.OkResp("success_tags_retrieved", tags))
}

func (h *Handler) CreateTag(c *fiber.Ctx) error {
	var input models.CreateTagInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	tag, err := h.tagService.Create(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_tag", err.Error()))
	}
	return c.JSON(web.OkResp("success_tag_created", tag))
}

func (h *Handler) UpdateTag(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid tag ID"))
	}

	var input models.CreateTagInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.tagService.Update(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_tag", err.Error()))
	}
	return c.JSON(web.OkResp("success_tag_updated", nil))
}

func (h *Handler) DeleteTag(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid tag ID"))
	}

	if err := h.tagService.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_tag", err.Error()))
	}
	return c.JSON(web.OkResp("success_tag_deleted", nil))
}
