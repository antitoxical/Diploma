package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

// AddToFavorites godoc
// @Summary Add to favorites
// @Description Adds a product to user's favorites
// @Tags favorites
// @Accept json
// @Produce json
// @Param input body models.CreateFavorite true "Favorite item"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /favorites [post]
func (h *Handler) AddToFavorites(c *fiber.Ctx) error {
	var input models.CreateFavorite
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.favoriteService.AddToFavorites(c.Context(), input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_add_to_favorites", err.Error()))
	}

	return c.JSON(web.OkResp("success_item_added_to_favorites", input))
}

// GetUserFavorites godoc
// @Summary Get user's favorites
// @Description Returns all user's favorite products
// @Tags favorites
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /favorites/{user_id} [get]
func (h *Handler) GetUserFavorites(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	favorites, err := h.favoriteService.GetUserFavorites(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_user_favorites", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_favorites_retrieved", favorites))
}

// RemoveFromFavorites godoc
// @Summary Remove from favorites
// @Description Removes a product from user's favorites
// @Tags favorites
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Param product_id path string true "Product UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /favorites/{user_id}/{product_id} [delete]
func (h *Handler) RemoveFromFavorites(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	productID := c.Params("product_id")

	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	if err := h.favoriteService.RemoveFromFavorites(c.Context(), telegramID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_remove_from_favorites", err.Error()))
	}

	return c.JSON(web.OkResp("success_item_removed_from_favorites", nil))
}
