package handler

import (
	"strconv"
	"strings"

	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *Handler) CreateProduct(c *fiber.Ctx) error {
	var input models.Product
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	product, err := h.productService.CreateProduct(c.Context(), input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_create_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_created", product))
}

func (h *Handler) GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}

	product, err := h.productService.GetProductByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_retrieved", product))
}

func (h *Handler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.productService.GetAllProducts(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_get_all_products", err.Error()))
	}

	return c.JSON(web.OkResp("success_products_retrieved", products))
}

func (h *Handler) SearchProducts(c *fiber.Ctx) error {
	var filter models.ProductFilter

	if v := c.Query("category_id"); v != "" {
		if _, err := uuid.Parse(v); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_category_id", "Invalid category_id"))
		}
		filter.CategoryID = &v
	}
	if c.Query("include_subcategories") == "true" {
		filter.IncludeSubcategories = true
	}
	if v := c.Query("firm_id"); v != "" {
		if _, err := uuid.Parse(v); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_firm_id", "Invalid firm_id"))
		}
		filter.FirmID = &v
	}
	if v := c.Query("tag_id"); v != "" {
		if _, err := uuid.Parse(v); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_tag_id", "Invalid tag_id"))
		}
		filter.TagID = &v
	}
	if v := c.Query("tag_ids"); v != "" {
		ids := strings.Split(v, ",")
		for _, id := range ids {
			id = strings.TrimSpace(id)
			if _, err := uuid.Parse(id); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_tag_ids", "Invalid UUID in tag_ids"))
			}
			filter.TagIDs = append(filter.TagIDs, id)
		}
	}
	if v := c.Query("search"); v != "" {
		filter.Search = &v
	}
	if v := c.Query("min_price"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_min_price", "min_price must be a non-negative integer (kopeks)"))
		}
		filter.MinPrice = &n
	}
	if v := c.Query("max_price"); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err != nil || n < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_max_price", "max_price must be a non-negative integer (kopeks)"))
		}
		filter.MaxPrice = &n
	}
	if v := c.Query("in_stock"); v == "true" {
		b := true
		filter.InStock = &b
	}
	if v := c.Query("min_rating"); v != "" {
		n, err := strconv.ParseFloat(v, 64)
		if err != nil || n < 0 || n > 5 {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_min_rating", "min_rating must be between 0 and 5"))
		}
		filter.MinRating = &n
	}
	if v := c.Query("new_days"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n <= 0 {
			return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_new_days", "new_days must be a positive integer"))
		}
		filter.NewDays = &n
	}
	if v := c.Query("has_reviews"); v == "true" {
		b := true
		filter.HasReviews = &b
	}
	if v := c.Query("sort_by"); v != "" {
		filter.SortBy = v
	}

	filter.Page, _ = strconv.Atoi(c.Query("page", "1"))
	filter.Limit, _ = strconv.Atoi(c.Query("limit", "20"))

	result, err := h.productService.SearchProducts(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_search_products", err.Error()))
	}

	return c.JSON(web.OkResp("success_products_search", result))
}

func (h *Handler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}
	var input models.UpdateProductInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_request_body", "Invalid request body"))
	}

	if err := h.productService.UpdateProduct(c.Context(), id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_updated", nil))
}

func (h *Handler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}

	if err := h.productService.DeleteProduct(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_delete_product", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_deleted", nil))
}

func (h *Handler) AddProductImage(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}

	var input models.ImageInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	img, err := h.productService.AddProductImage(c.Context(), id, input.Image, false)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_add_product_image", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_image_added", img))
}

func (h *Handler) RemoveProductImage(c *fiber.Ctx) error {
	imageID := c.Params("image_id")
	if _, err := uuid.Parse(imageID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid image ID"))
	}

	if err := h.productService.RemoveProductImage(c.Context(), imageID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_remove_product_image", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_image_removed", nil))
}

func (h *Handler) SetProductImages(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}

	var input models.ImagesInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.productService.SetProductImages(c.Context(), id, input.Images); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_set_product_images", err.Error()))
	}

	return c.JSON(web.OkResp("success_product_images_set", nil))
}

func (h *Handler) UpdateStock(c *fiber.Ctx) error {
	id := c.Params("id")
	if _, err := uuid.Parse(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}

	var input models.StockInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_body", "Invalid body"))
	}

	if err := h.productService.UpdateStock(c.Context(), id, input.Stock); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_update_stock", err.Error()))
	}

	return c.JSON(web.OkResp("success_stock_updated", nil))
}

func (h *Handler) AddProductTag(c *fiber.Ctx) error {
	productID := c.Params("id")
	tagID := c.Params("tag_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}
	if _, err := uuid.Parse(tagID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_tag_id", "Invalid tag ID"))
	}

	if err := h.productService.AddProductTag(c.Context(), productID, tagID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_add_product_tag", err.Error()))
	}
	return c.JSON(web.OkResp("success_product_tag_added", nil))
}

func (h *Handler) RemoveProductTag(c *fiber.Ctx) error {
	productID := c.Params("id")
	tagID := c.Params("tag_id")
	if _, err := uuid.Parse(productID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_id", "Invalid product ID"))
	}
	if _, err := uuid.Parse(tagID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(web.ErrorResp("error_invalid_tag_id", "Invalid tag ID"))
	}

	if err := h.productService.RemoveProductTag(c.Context(), productID, tagID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(web.ErrorResp("error_remove_product_tag", err.Error()))
	}
	return c.JSON(web.OkResp("success_product_tag_removed", nil))
}
