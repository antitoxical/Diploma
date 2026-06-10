package models

type OrderByDay struct {
	Date          string `db:"date" json:"date"`
	Count         int    `db:"count" json:"count"`
	RevenueKopeks int64  `db:"revenue_kopeks" json:"revenue_kopeks"`
}

type OrderByStatus struct {
	Status string `db:"status" json:"status"`
	Count  int    `db:"count" json:"count"`
}

type TopProduct struct {
	ID            string `db:"id" json:"id"`
	Name          string `db:"name" json:"name"`
	SoldCount     int    `db:"sold_count" json:"sold_count"`
	RevenueKopeks int64  `db:"revenue_kopeks" json:"revenue_kopeks"`
}

type DashboardStats struct {
	TotalOrders    int64           `json:"total_orders"`
	TotalUsers     int64           `json:"total_users"`
	TotalRevenue   int64           `json:"total_revenue_kopeks"`
	TotalProducts  int64           `json:"total_products"`
	OrdersByDay    []OrderByDay    `json:"orders_by_day"`
	OrdersByStatus []OrderByStatus `json:"orders_by_status"`
	TopProducts    []TopProduct    `json:"top_products"`
}
