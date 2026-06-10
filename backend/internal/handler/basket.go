package handler

import (
	"strconv"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

// AddToBasket godoc
// @Summary Add item to basket
// @Description Adds a product to user's basket
// @Tags basket
// @Accept json
// @Produce json
// @Param input body models.CreateBasketItem true "Basket item"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /basket [post]
func (h *Handler) AddToBasket(c *fiber.Ctx) error {
	var input models.CreateBasketItem
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.basketService.AddToBasket(c.Context(), input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_add_to_basket", err.Error()))
	}

	return c.JSON(web.OkResp("success_item_added_to_basket", input))
}

// GetUserBasket godoc
// @Summary Get user's basket
// @Description Returns all items in user's basket
// @Tags basket
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /basket/{user_id} [get]
func (h *Handler) GetUserBasket(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	items, err := h.basketService.GetUserBasket(c.Context(), telegramID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_user_basket", err.Error()))
	}

	return c.JSON(web.OkResp("success_user_basket_retrieved", items))
}

// UpdateBasketItem godoc
// @Summary Update basket item
// @Description Updates quantity of product in basket
// @Tags basket
// @Accept json
// @Produce json
// @Param input body models.CreateBasketItem true "Updated basket item"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /basket [put]
func (h *Handler) UpdateBasketItem(c *fiber.Ctx) error {
	var input models.CreateBasketItem
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.basketService.UpdateBasketItem(c.Context(), input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_basket_item", err.Error()))
	}

	return c.JSON(web.OkResp("success_basket_item_updated", input))
}

// RemoveFromBasket godoc
// @Summary Remove item from basket
// @Description Removes a product from user's basket
// @Tags basket
// @Produce json
// @Param user_id path integer true "Telegram ID"
// @Param product_id path string true "Product UUID"
// @Success 200 {object} web.Response
// @Failure 400 {object} web.Response
// @Failure 500 {object} web.Response
// @Router /basket/{user_id}/{product_id} [delete]
func (h *Handler) RemoveFromBasket(c *fiber.Ctx) error {
	userIDStr := c.Params("user_id")
	productID := c.Params("product_id")

	telegramID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_user_id", "Invalid user ID"))
	}

	if err := h.basketService.RemoveFromBasket(c.Context(), telegramID, productID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_remove_from_basket", err.Error()))
	}

	return c.JSON(web.OkResp("success_item_removed_from_basket", nil))
}
