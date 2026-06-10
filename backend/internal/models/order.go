package models

import "time"

type OrderItemRow struct {
	ID                    string `db:"id" json:"id"`
	OrderID               string `db:"order_id" json:"order_id"`
	ProductID             string `db:"product_id" json:"product_id"`
	ProductName           string `db:"product_name" json:"product_name,omitempty"`
	Quantity              int    `db:"quantity" json:"quantity"`
	PriceAtPurchaseKopeks int64  `db:"price_at_purchase_kopeks" json:"price_at_purchase_kopeks"`
}

type OrderWithProducts struct {
	ID                 string         `db:"id" json:"id"`
	UserID             string         `db:"user_id" json:"user_id"`
	TelegramID         int64          `db:"telegram_id" json:"telegram_id,omitempty"`
	DeliveryMethodID   *string        `db:"delivery_method_id" json:"delivery_method_id,omitempty"`
	PaymentMethodID    *string        `db:"payment_method_id" json:"payment_method_id,omitempty"`
	DeliveryMethodName *string        `db:"delivery_method_name" json:"delivery_method_name,omitempty"`
	PaymentMethodName  *string        `db:"payment_method_name" json:"payment_method_name,omitempty"`
	PromocodeID        *string        `db:"promocode_id" json:"promocode_id,omitempty"`
	Status             string         `db:"status" json:"status"`
	TotalAmountKopeks  int64          `db:"total_amount_kopeks" json:"total_amount_kopeks"`
	AddressSnapshot    *string        `db:"address_snapshot" json:"address_snapshot,omitempty"`
	CreatedAt          time.Time      `db:"created_at" json:"created_at"`
	Products           []OrderItemRow `json:"products"`
}

type OrderStatusHistory struct {
	ID        string    `db:"id" json:"id"`
	OrderID   string    `db:"order_id" json:"order_id"`
	Status    string    `db:"status" json:"status"`
	Comment   *string   `db:"comment" json:"comment,omitempty"`
	ChangedAt time.Time `db:"changed_at" json:"changed_at"`
}

type CreateOrder struct {
	UserID           int64  `json:"user_id"`
	DeliveryMethodID string `json:"delivery_method_id"`
	PaymentMethodID  string `json:"payment_method_id"`
	PromoCode        string `json:"promo_code,omitempty"`
	DeliveryAddress  string `json:"delivery_address"`
	Items            []struct {
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	} `json:"items"`
}

type UpdateOrderStatus struct {
	Status  string `json:"status"`
	Comment string `json:"comment,omitempty"`
}

type ProductRatingStat struct {
	ProductID   string  `db:"product_id" json:"product_id"`
	AvgRating   float64 `db:"avg_rating" json:"avg_rating"`
	ReviewCount int     `db:"review_count" json:"review_count"`
}

type DeliveryMethod struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Kind        string    `db:"kind" json:"kind"`
	PriceKopeks int64     `db:"price_kopeks" json:"price_kopeks"`
	Description *string   `db:"description" json:"description,omitempty"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type PaymentMethod struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Description *string   `db:"description" json:"description,omitempty"`
	IsActive    bool      `db:"is_active" json:"is_active"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}

type Promocode struct {
	ID              string     `db:"id" json:"id"`
	Code            string     `db:"code" json:"code"`
	DiscountPercent int        `db:"discount_percent" json:"discount_percent"`
	MaxUses         *int       `db:"max_uses" json:"max_uses,omitempty"`
	CurrentUses     int        `db:"current_uses" json:"current_uses"`
	IsActive        bool       `db:"is_active" json:"is_active"`
	ExpiresAt       *time.Time `db:"expires_at" json:"expires_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
}

type CreatePromocode struct {
	Code            string `json:"code"`
	DiscountPercent int    `json:"discount_percent"`
	MaxUses         *int   `json:"max_uses,omitempty"`
	ExpiresAt       string `json:"expires_at,omitempty"`
}
