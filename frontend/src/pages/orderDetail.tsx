import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api/client";
import { Order, OrderStatusHistory, formatPrice, orderStatusLabel, orderStatusColor } from "../api/types";
import { formatOrderAddress, orderItemTitle, orderTotalQuantity } from "../utils/orderDisplay";
import "../style/orders.css";

const OrderDetail: React.FC = () => {
  const { orderId } = useParams();
  const navigate = useNavigate();
  const [order, setOrder] = useState<Order | null>(null);
  const [history, setHistory] = useState<OrderStatusHistory[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!orderId) return;
    setLoading(true);
    Promise.all([
      api.getOrder(orderId),
      api.getOrderStatusHistory(orderId).catch(() => []),
    ])
      .then(([o, h]) => { setOrder(o); setHistory(Array.isArray(h) ? h : []); })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [orderId]);

  if (loading) {
    return (
      <div className="orders-container">
        <p className="orders-hint">Загрузка...</p>
      </div>
    );
  }

  if (!order) {
    return (
      <div className="orders-container">
        <p className="orders-hint">Заказ не найден</p>
      </div>
    );
  }

  const addressLines = formatOrderAddress(order.address_snapshot);
  const totalQty = orderTotalQuantity(order.products);

  return (
    <div className="orders-container">
      <div className="orders-header">
        <button className="orders-back" onClick={() => navigate(-1)}>&lt;</button>
        <h1 className="orders-title">Заказ #{order.id.slice(0, 8)}</h1>
      </div>

      <div className="order-detail-card">
        <div className="order-card-top">
          <span
            className="order-card-status"
            style={{ background: orderStatusColor(order.status) }}
          >
            {orderStatusLabel(order.status)}
          </span>
          <span className="order-card-date">
            {new Date(order.created_at).toLocaleDateString("ru-RU", {
              day: "numeric", month: "long", year: "numeric", hour: "2-digit", minute: "2-digit",
            })}
          </span>
        </div>

        <div className="order-detail-meta">
          <div className="order-detail-meta-row">
            <span className="order-detail-meta-label">Позиций</span>
            <span>{order.products?.length ?? 0}</span>
          </div>
          <div className="order-detail-meta-row">
            <span className="order-detail-meta-label">Единиц товара</span>
            <span>{totalQty}</span>
          </div>
          {order.delivery_method_name && (
            <div className="order-detail-meta-row">
              <span className="order-detail-meta-label">Доставка</span>
              <span>{order.delivery_method_name}</span>
            </div>
          )}
          {order.payment_method_name && (
            <div className="order-detail-meta-row">
              <span className="order-detail-meta-label">Оплата</span>
              <span>{order.payment_method_name}</span>
            </div>
          )}
        </div>

        {addressLines.length > 0 && (
          <div className="order-detail-address">
            <h3>Адрес и доставка</h3>
            {addressLines.map((line, idx) => (
              <p key={idx}>{line}</p>
            ))}
          </div>
        )}

        <div className="order-detail-items">
          <h3>Товары</h3>
          {(order.products || []).map(item => (
            <div key={item.id} className="order-detail-item-row">
              <span className="order-detail-item-name">
                {orderItemTitle(item)} × {item.quantity}
              </span>
              <span>{formatPrice(item.price_at_purchase_kopeks * item.quantity)}</span>
            </div>
          ))}
        </div>

        <div className="order-card-total order-detail-total">
          Итого: {formatPrice(order.total_amount_kopeks)}
        </div>
      </div>

      {history.length > 0 && (
        <div className="order-history">
          <h3>История статусов</h3>
          <div className="order-history-list">
            {history.map(h => (
              <div key={h.id} className="order-history-item">
                <span
                  className="order-history-status"
                  style={{ color: orderStatusColor(h.status) }}
                >
                  {orderStatusLabel(h.status)}
                </span>
                <span className="order-history-date">
                  {new Date(h.changed_at).toLocaleString("ru-RU")}
                </span>
                {h.comment && <p className="order-history-comment">{h.comment}</p>}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default OrderDetail;
