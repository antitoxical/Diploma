package handler

import (
	"database/sql"
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) GetAllReviews(c *fiber.Ctx) error {
	list, err := h.reviewService.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_reviews", err.Error()))
	}
	if list == nil {
		list = []models.Review{}
	}
	return c.JSON(web.OkResp("success_reviews_retrieved", list))
}

func (h *Handler) UpsertReview(c *fiber.Ctx) error {
	var input models.UpsertReviewInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}
	if _, err := uuid.Parse(input.ProductID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	rev, err := h.reviewService.UpsertReview(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_upsert_review", err.Error()))
	}
	return c.JSON(web.OkResp("success_review_upserted", rev))
}

func (h *Handler) DeleteReview(c *fiber.Ctx) error {
	telegramID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user_id"))
	}
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	if err := h.reviewService.DeleteReview(c.Context(), telegramID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_review", err.Error()))
	}
	return c.JSON(web.OkResp("success_review_deleted", nil))
}

func (h *Handler) GetReviewsByProduct(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	list, err := h.reviewService.GetByProduct(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_reviews", err.Error()))
	}
	return c.JSON(web.OkResp("success_reviews_retrieved", list))
}

func (h *Handler) GetUserReview(c *fiber.Ctx) error {
	telegramID, err := strconv.ParseInt(c.Params("user_id"), 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user_id"))
	}
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	rev, err := h.reviewService.GetUserReview(c.Context(), telegramID, productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(web.ErrorResp("error_review_not_found", "Review not found"))
		}
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_review", err.Error()))
	}
	return c.JSON(web.OkResp("success_review_retrieved", rev))
}

func (h *Handler) GetProductRatingStat(c *fiber.Ctx) error {
	productID := c.Params("product_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_product_id", "Invalid product ID"))
	}

	stat, err := h.reviewService.GetRatingStat(c.Context(), productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_rating_stat", err.Error()))
	}
	return c.JSON(web.OkResp("success_rating_stat", stat))
}

func (h *Handler) ModerateReview(c *fiber.Ctx) error {
	reviewID := c.Params("id")
	if _, err := uuid.Parse(reviewID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid review ID"))
	}

	var input struct {
		Approved bool `json:"approved"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.reviewService.ModerateReview(c.Context(), reviewID, input.Approved); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_moderate_review", err.Error()))
	}
	return c.JSON(web.OkResp("success_review_moderated", nil))
}
