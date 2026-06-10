package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	var input models.Category
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	category, err := h.categoryService.CreateCategory(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_category", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_created", category))
}

func (h *Handler) GetCategoryByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}

	category, err := h.categoryService.GetCategoryByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_category", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_retrieved", category))
}

func (h *Handler) GetAllCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetAllCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_all_categories", err.Error()))
	}

	return c.JSON(web.OkResp("success_categories_retrieved", categories))
}

func (h *Handler) GetRootCategories(c *fiber.Ctx) error {
	categories, err := h.categoryService.GetRootCategories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_root_categories", err.Error()))
	}

	return c.JSON(web.OkResp("success_root_categories_retrieved", categories))
}

func (h *Handler) GetSubcategories(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}

	categories, err := h.categoryService.GetSubcategories(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_subcategories", err.Error()))
	}

	return c.JSON(web.OkResp("success_subcategories_retrieved", categories))
}

func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}
	var input models.UpdateCategoryInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.categoryService.UpdateCategory(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_category", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_updated", nil))
}

func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}

	if err := h.categoryService.DeleteCategory(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_category", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_deleted", nil))
}

func (h *Handler) SetCategoryImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}

	var input models.ImageInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.categoryService.SetImage(c.Context(), id, input.Image); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_set_image", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_image_set", nil))
}

func (h *Handler) RemoveCategoryImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid category ID"))
	}

	if err := h.categoryService.RemoveImage(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_remove_image", err.Error()))
	}

	return c.JSON(web.OkResp("success_category_image_removed", nil))
}
