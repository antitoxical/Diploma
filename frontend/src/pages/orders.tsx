import React, { useEffect, useState } from "react";
import { useTelegram } from "../context/TelegramContext";
import { api } from "../api/client";
import { Order, formatPrice, orderStatusLabel, orderStatusColor } from "../api/types";
import { orderProductsPreview, orderTotalQuantity } from "../utils/orderDisplay";
import { useNavigate } from "react-router-dom";
import "../style/orders.css";

const Orders: React.FC = () => {
  const { user: tgUser } = useTelegram();
  const [orders, setOrders] = useState<Order[]>([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    if (!tgUser) { setLoading(false); return; }
    api.getUserOrders(tgUser.id)
      .then(data => setOrders(Array.isArray(data) ? data : []))
      .catch(() => setOrders([]))
      .finally(() => setLoading(false));
  }, [tgUser]);

  return (
    <div className="orders-container">
      <div className="orders-header">
        <button className="orders-back" onClick={() => navigate(-1)}>&lt;</button>
        <h1 className="orders-title">Мои заказы</h1>
      </div>

      {loading && <p className="orders-hint">Загрузка...</p>}

      {!loading && orders.length === 0 && (
        <div className="orders-empty">
          <p>У вас пока нет заказов</p>
          <button className="orders-shop-btn" onClick={() => navigate("/catalog")}>
            Перейти в каталог
          </button>
        </div>
      )}

      <div className="orders-list">
        {orders.map(order => {
          const totalQty = orderTotalQuantity(order.products || []);
          return (
            <div
              className="order-card"
              key={order.id}
              onClick={() => navigate(`/order/${order.id}`)}
              style={{ cursor: "pointer" }}
            >
              <div className="order-card-top">
                <span className="order-card-id">#{order.id.slice(0, 8)}</span>
                <span
                  className="order-card-status"
                  style={{ background: orderStatusColor(order.status) }}
                >
                  {orderStatusLabel(order.status)}
                </span>
              </div>
              <div className="order-card-date">
                {new Date(order.created_at).toLocaleDateString("ru-RU", {
                  day: "numeric", month: "long", year: "numeric",
                })}
              </div>
              <div className="order-card-items">
                {orderProductsPreview(order.products || [])}
                {totalQty > 0 && (
                  <span className="order-card-qty"> ({totalQty} шт.)</span>
                )}
              </div>
              {(order.delivery_method_name || order.payment_method_name) && (
                <div className="order-card-methods">
                  {order.delivery_method_name && <span>{order.delivery_method_name}</span>}
                  {order.delivery_method_name && order.payment_method_name && " · "}
                  {order.payment_method_name && <span>{order.payment_method_name}</span>}
                </div>
              )}
              <div className="order-card-total">
                {formatPrice(order.total_amount_kopeks)}
              </div>
              {order.address_snapshot && (
                <div className="order-card-address">
                  {order.address_snapshot.split("\n")[0]}
                </div>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
};

export default Orders;
