package orders

import (
	"context"
	"database/sql"

	"telegramshop_backend/internal/models"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateOrder(ctx context.Context, input models.CreateOrder, userUUID string, totalKopeks int64, lines []OrderLineInput, deliveryMethodID, paymentMethodID, promocodeID *string) (models.OrderWithProducts, error)
	GetOrderByID(ctx context.Context, id string) (models.OrderWithProducts, error)
	GetUserOrders(ctx context.Context, telegramID int64) ([]models.OrderWithProducts, error)
	GetAll(ctx context.Context) ([]models.OrderWithProducts, error)
	UpdateStatus(ctx context.Context, id string, status string, comment *string) error
	GetStatusHistory(ctx context.Context, orderID string) ([]models.OrderStatusHistory, error)
	LookupMethodNames(ctx context.Context, deliveryMethodID, paymentMethodID *string) (deliveryName, paymentName string)
}

type OrderLineInput struct {
	ProductID             string
	Quantity              int
	PriceAtPurchaseKopeks int64
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateOrder(ctx context.Context, input models.CreateOrder, userUUID string, totalKopeks int64, lines []OrderLineInput, deliveryMethodID, paymentMethodID, promocodeID *string) (models.OrderWithProducts, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return models.OrderWithProducts{}, err
	}
	defer tx.Rollback()

	orderQuery := `
		INSERT INTO orders (user_id, delivery_method_id, payment_method_id, promocode_id, status, total_amount_kopeks, address_snapshot, created_at)
		VALUES ($1, $2, $3, $4, 'created', $5, $6, now())
		RETURNING id, user_id, delivery_method_id, payment_method_id, promocode_id, status, total_amount_kopeks, address_snapshot, created_at`

	var order models.OrderWithProducts
	var addr, dmID, pmID, pcID sql.NullString
	err = tx.QueryRowContext(ctx, orderQuery, userUUID, deliveryMethodID, paymentMethodID, promocodeID, totalKopeks, input.DeliveryAddress).Scan(
		&order.ID, &order.UserID, &dmID, &pmID, &pcID, &order.Status, &order.TotalAmountKopeks, &addr, &order.CreatedAt,
	)
	if err != nil {
		return models.OrderWithProducts{}, err
	}
	if addr.Valid {
		s := addr.String
		order.AddressSnapshot = &s
	}
	if dmID.Valid {
		s := dmID.String
		order.DeliveryMethodID = &s
	}
	if pmID.Valid {
		s := pmID.String
		order.PaymentMethodID = &s
	}
	if pcID.Valid {
		s := pcID.String
		order.PromocodeID = &s
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO order_status_history (order_id, status, comment, changed_at)
		VALUES ($1, 'created', 'Заказ создан', now())`, order.ID)
	if err != nil {
		return models.OrderWithProducts{}, err
	}

	itemQuery := `
		INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase_kopeks)
		VALUES ($1, $2, $3, $4)
		RETURNING id, order_id, product_id, quantity, price_at_purchase_kopeks`

	for _, line := range lines {
		var row models.OrderItemRow
		err = tx.QueryRowContext(ctx, itemQuery, order.ID, line.ProductID, line.Quantity, line.PriceAtPurchaseKopeks).Scan(
			&row.ID, &row.OrderID, &row.ProductID, &row.Quantity, &row.PriceAtPurchaseKopeks,
		)
		if err != nil {
			return models.OrderWithProducts{}, err
		}
		order.Products = append(order.Products, row)
	}

	for _, line := range lines {
		res, err := tx.ExecContext(ctx, `
			UPDATE products SET stock_quantity = stock_quantity - $2
			WHERE id = $1::uuid AND stock_quantity >= $2`, line.ProductID, line.Quantity)
		if err != nil {
			return models.OrderWithProducts{}, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return models.OrderWithProducts{}, err
		}
		if n == 0 {
			return models.OrderWithProducts{}, sql.ErrNoRows
		}
	}

	if err = tx.Commit(); err != nil {
		return models.OrderWithProducts{}, err
	}

	return order, nil
}

func (r *repository) GetOrderByID(ctx context.Context, id string) (models.OrderWithProducts, error) {
	orderQuery := `
		SELECT o.id, o.user_id, u.telegram_id, o.delivery_method_id, o.payment_method_id, o.promocode_id,
		       o.status, o.total_amount_kopeks, o.address_snapshot, o.created_at,
		       dm.name AS delivery_method_name, pm.name AS payment_method_name
		FROM orders o
		JOIN users u ON u.id = o.user_id
		LEFT JOIN delivery_methods dm ON dm.id = o.delivery_method_id
		LEFT JOIN payment_methods pm ON pm.id = o.payment_method_id
		WHERE o.id = $1`

	var order models.OrderWithProducts
	var addr, dmID, pmID, pcID, dmName, pmName sql.NullString
	err := r.db.QueryRowContext(ctx, orderQuery, id).Scan(
		&order.ID, &order.UserID, &order.TelegramID, &dmID, &pmID, &pcID,
		&order.Status, &order.TotalAmountKopeks, &addr, &order.CreatedAt,
		&dmName, &pmName,
	)
	if err != nil {
		return models.OrderWithProducts{}, err
	}
	if addr.Valid {
		s := addr.String
		order.AddressSnapshot = &s
	}
	if dmID.Valid {
		s := dmID.String
		order.DeliveryMethodID = &s
	}
	if pmID.Valid {
		s := pmID.String
		order.PaymentMethodID = &s
	}
	if pcID.Valid {
		s := pcID.String
		order.PromocodeID = &s
	}
	if dmName.Valid {
		s := dmName.String
		order.DeliveryMethodName = &s
	}
	if pmName.Valid {
		s := pmName.String
		order.PaymentMethodName = &s
	}

	err = r.db.SelectContext(ctx, &order.Products, `
		SELECT oi.id, oi.order_id, oi.product_id, COALESCE(p.name, '') AS product_name,
		       oi.quantity, oi.price_at_purchase_kopeks
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1`, id)
	if err != nil {
		return models.OrderWithProducts{}, err
	}
	return order, nil
}

func (r *repository) GetUserOrders(ctx context.Context, telegramID int64) ([]models.OrderWithProducts, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids, `
		SELECT o.id FROM orders o
		JOIN users u ON u.id = o.user_id
		WHERE u.telegram_id = $1
		ORDER BY o.created_at DESC`, telegramID)
	if err != nil {
		return nil, err
	}
	out := make([]models.OrderWithProducts, 0, len(ids))
	for _, id := range ids {
		o, err := r.GetOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

func (r *repository) GetAll(ctx context.Context) ([]models.OrderWithProducts, error) {
	var ids []string
	err := r.db.SelectContext(ctx, &ids, `SELECT id FROM orders ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	out := make([]models.OrderWithProducts, 0, len(ids))
	for _, id := range ids {
		o, err := r.GetOrderByID(ctx, id)
		if err != nil {
			return nil, err
		}
		out = append(out, o)
	}
	return out, nil
}

func (r *repository) UpdateStatus(ctx context.Context, id string, status string, comment *string) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	res, err := tx.ExecContext(ctx, `UPDATE orders SET status = $1 WHERE id = $2::uuid`, status, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO order_status_history (order_id, status, comment, changed_at)
		VALUES ($1, $2, $3, now())`, id, status, comment)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) GetStatusHistory(ctx context.Context, orderID string) ([]models.OrderStatusHistory, error) {
	var history []models.OrderStatusHistory
	err := r.db.SelectContext(ctx, &history, `
		SELECT id, order_id, status, comment, changed_at
		FROM order_status_history
		WHERE order_id = $1
		ORDER BY changed_at`, orderID)
	return history, err
}

func (r *repository) LookupMethodNames(ctx context.Context, deliveryMethodID, paymentMethodID *string) (string, string) {
	var deliveryName, paymentName string
	if deliveryMethodID != nil && *deliveryMethodID != "" {
		_ = r.db.GetContext(ctx, &deliveryName, `SELECT name FROM delivery_methods WHERE id = $1`, *deliveryMethodID)
	}
	if paymentMethodID != nil && *paymentMethodID != "" {
		_ = r.db.GetContext(ctx, &paymentName, `SELECT name FROM payment_methods WHERE id = $1`, *paymentMethodID)
	}
	return deliveryName, paymentName
}
