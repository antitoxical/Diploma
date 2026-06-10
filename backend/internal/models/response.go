package models

// Response structures for Swagger documentation

// BasketResponse represents a basket item response
type BasketResponse struct {
	Status string     `json:"status" example:"success_item_added_to_basket"`
	Data   BasketItem `json:"data"`
}

// BasketListResponse represents a list of basket items response
type BasketListResponse struct {
	Status string       `json:"status" example:"success_user_basket_retrieved"`
	Data   []BasketItem `json:"data"`
}

// FavoriteResponse represents a favorite item response
type FavoriteResponse struct {
	Status string   `json:"status" example:"success_item_added_to_favorites"`
	Data   Favorite `json:"data"`
}

// FavoriteListResponse represents a list of favorite items response
type FavoriteListResponse struct {
	Status string     `json:"status" example:"success_user_favorites_retrieved"`
	Data   []Favorite `json:"data"`
}

// OrderResponse represents an order response
type OrderResponse struct {
	Status string            `json:"status" example:"success_order_created"`
	Data   OrderWithProducts `json:"data"`
}

// OrderListResponse represents a list of orders response
type OrderListResponse struct {
	Status string              `json:"status" example:"success_all_orders_retrieved"`
	Data   []OrderWithProducts `json:"data"`
}

// UserResponse represents a user response
type UserResponse struct {
	Status string `json:"status" example:"success_user_created"`
	Data   User   `json:"data"`
}

// UserListResponse represents a list of users response
type UserListResponse struct {
	Status string `json:"status" example:"success_all_users_retrieved"`
	Data   []User `json:"data"`
}

// ProductResponse represents a product response
type ProductResponse struct {
	Status string  `json:"status" example:"success_product_created"`
	Data   Product `json:"data"`
}

// ProductListResponse represents a list of products response
type ProductListResponse struct {
	Status string    `json:"status" example:"success_all_products_retrieved"`
	Data   []Product `json:"data"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	Status string   `json:"status" example:"success_category_created"`
	Data   Category `json:"data"`
}

// CategoryListResponse represents a list of categories response
type CategoryListResponse struct {
	Status string     `json:"status" example:"success_all_categories_retrieved"`
	Data   []Category `json:"data"`
}

// FirmResponse represents a firm response
type FirmResponse struct {
	Status string `json:"status" example:"success_firm_created"`
	Data   Firm   `json:"data"`
}

// FirmListResponse represents a list of firms response
type FirmListResponse struct {
	Status string `json:"status" example:"success_all_firms_retrieved"`
	Data   []Firm `json:"data"`
}

// PriceResponse represents a product price tier response
type PriceResponse struct {
	Status string           `json:"status" example:"success_price_created"`
	Data   ProductPriceTier `json:"data"`
}

// PriceListResponse represents a list of price tiers response
type PriceListResponse struct {
	Status string             `json:"status" example:"success_prices_retrieved"`
	Data   []ProductPriceTier `json:"data"`
}

// NotificationResponse represents a notification response
type NotificationResponse struct {
	Status string       `json:"status" example:"success_notification_created"`
	Data   Notification `json:"data"`
}

// NotificationListResponse represents a list of notifications response
type NotificationListResponse struct {
	Status string         `json:"status" example:"success_notifications_retrieved"`
	Data   []Notification `json:"data"`
}

// SuccessResponse represents a generic success response
type SuccessResponse struct {
	Status string      `json:"status" example:"success_operation_completed"`
	Data   interface{} `json:"data"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Status string `json:"status" example:"error_invalid_request_body"`
	Data   string `json:"data" example:"Invalid request body"`
}
