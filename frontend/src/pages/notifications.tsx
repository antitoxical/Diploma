import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTelegram } from "../context/TelegramContext";
import { api } from "../api/client";
import { Notification as AppNotification } from "../api/types";
import "../style/orders.css";

const Notifications: React.FC = () => {
  const { user: tgUser } = useTelegram();
  const navigate = useNavigate();
  const [items, setItems] = useState<AppNotification[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!tgUser) { setLoading(false); return; }
    api.getNotifications(tgUser.id)
      .then(data => setItems(Array.isArray(data) ? data : []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, [tgUser]);

  const handleMarkAll = () => {
    if (!tgUser) return;
    api.markAllNotificationsRead(tgUser.id)
      .then(() => setItems(prev => prev.map(n => ({ ...n, is_read: true }))))
      .catch(() => {});
  };

  const handleRead = (id: string) => {
    api.markNotificationRead(id)
      .then(() => setItems(prev => prev.map(n => n.id === id ? { ...n, is_read: true } : n)))
      .catch(() => {});
  };

  return (
    <div className="orders-container">
      <div className="orders-header">
        <button className="orders-back" onClick={() => navigate(-1)}>&lt;</button>
        <h1 className="orders-title">Уведомления</h1>
      </div>

      {items.some(n => !n.is_read) && (
        <button className="orders-shop-btn" onClick={handleMarkAll} style={{ marginBottom: 12 }}>
          Прочитать все
        </button>
      )}

      {loading && <p className="orders-hint">Загрузка...</p>}

      {!loading && items.length === 0 && (
        <div className="orders-empty">
          <p>Уведомлений пока нет</p>
        </div>
      )}

      <div className="orders-list">
        {items.map(n => (
          <div
            key={n.id}
            className="order-card"
            style={{ opacity: n.is_read ? 0.6 : 1, cursor: "pointer" }}
            onClick={() => !n.is_read && handleRead(n.id)}
          >
            <div className="order-card-top">
              <span className="order-card-id">{n.title || "Уведомление"}</span>
              {!n.is_read && <span className="order-card-status" style={{ background: "#2196f3" }}>Новое</span>}
            </div>
            <div className="order-card-date">
              {new Date(n.created_at).toLocaleString("ru-RU")}
            </div>
            {n.message && <div className="order-card-address">{n.message}</div>}
          </div>
        ))}
      </div>
    </div>
  );
};

export default Notifications;
