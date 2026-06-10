package handler

import (
	"telegramshop_backend/internal/models"
	"telegramshop_backend/pkg/web"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) GetDashboardStats(c *fiber.Ctx) error {
	ctx := c.Context()
	var stats models.DashboardStats

	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM orders`).Scan(&stats.TotalOrders)
	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&stats.TotalUsers)
	_ = h.db.QueryRowContext(ctx,
		`SELECT COALESCE(SUM(total_amount_kopeks),0) FROM orders WHERE status = 'completed'`,
	).Scan(&stats.TotalRevenue)
	_ = h.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM products WHERE is_active = true`).Scan(&stats.TotalProducts)

	stats.OrdersByDay = []models.OrderByDay{}
	rows, err := h.db.QueryContext(ctx, `
		SELECT TO_CHAR(DATE_TRUNC('day', created_at AT TIME ZONE 'UTC'), 'YYYY-MM-DD') AS date,
		       COUNT(*)::int                                                              AS count,
		       COALESCE(SUM(total_amount_kopeks), 0)                                      AS revenue_kopeks
		FROM orders
		WHERE created_at >= NOW() - INTERVAL '30 days'
		GROUP BY DATE_TRUNC('day', created_at AT TIME ZONE 'UTC')
		ORDER BY 1`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var r models.OrderByDay
			if rows.Scan(&r.Date, &r.Count, &r.RevenueKopeks) == nil {
				stats.OrdersByDay = append(stats.OrdersByDay, r)
			}
		}
	}

	stats.OrdersByStatus = []models.OrderByStatus{}
	rows2, err := h.db.QueryContext(ctx, `
		SELECT status, COUNT(*)::int AS count FROM orders GROUP BY status`)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var r models.OrderByStatus
			if rows2.Scan(&r.Status, &r.Count) == nil {
				stats.OrdersByStatus = append(stats.OrdersByStatus, r)
			}
		}
	}

	stats.TopProducts = []models.TopProduct{}
	rows3, err := h.db.QueryContext(ctx, `
		SELECT p.id, p.name,
		       COALESCE(SUM(oi.quantity), 0)::int                              AS sold_count,
		       COALESCE(SUM(oi.quantity * oi.price_at_purchase_kopeks), 0)     AS revenue_kopeks
		FROM products p
		LEFT JOIN order_items oi ON oi.product_id = p.id
		GROUP BY p.id, p.name
		ORDER BY sold_count DESC
		LIMIT 10`)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var r models.TopProduct
			if rows3.Scan(&r.ID, &r.Name, &r.SoldCount, &r.RevenueKopeks) == nil {
				stats.TopProducts = append(stats.TopProducts, r)
			}
		}
	}

	return c.JSON(web.OkResp("success_dashboard_stats", stats))
}
