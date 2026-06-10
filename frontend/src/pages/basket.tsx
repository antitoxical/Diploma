import React, { useEffect } from "react";
import "../style/basket.css";
import { useCart } from "../context/CartContext";
import { useToast } from "../context/ToastContext";
import { Link } from "react-router-dom";
import { formatPrice, PLACEHOLDER_IMAGE } from "../api/types";

const Cart: React.FC = () => {
  const { cartItems, increment, decrement, removeFromCart, totalPrice, refreshStock } = useCart();
  const { showToast } = useToast();
  const isEmpty = cartItems.length === 0;

  useEffect(() => {
    if (cartItems.length > 0) {
      refreshStock();
    }
  }, [cartItems.length, refreshStock]);

  return (
    <div className="cart-container">
      <h1 className="h-home">Корзина</h1>
      {isEmpty ? (
        <div className="cart-empty">
          <div className="cart-empty-icon">
            <img src="/images/cart.svg" alt="cart" />
          </div>
          <div className="cart-empty-title">Ваша корзина пуста</div>
          <div className="cart-empty-desc">
            <span><Link to="/" className="cart-empty-link">Нажмите здесь</Link>, чтобы продолжить покупки</span>
          </div>
        </div>
      ) : (
        <div className="cart-full">
          {cartItems.map(item => (
            <div className="cart-item" key={item.id}>
              <img src={item.image || PLACEHOLDER_IMAGE} alt={item.title} className="cart-item-img" />
              <div className="cart-item-info">
                <div className="cart-item-title">{item.title}</div>
                <div className="cart-item-controls">
                  <button type="button" onClick={() => decrement(item.id)}>−</button>
                  <span>{item.count}</span>
                  <button
                    type="button"
                    onClick={() => {
                      if (!increment(item.id)) {
                        showToast(`Доступно не более ${item.stockMax} шт.`, "error");
                      }
                    }}
                    disabled={item.count >= item.stockMax}
                  >
                    +
                  </button>
                </div>
                {item.count >= item.stockMax && (
                  <div className="cart-item-stock-hint cart-item-stock-hint-limit">
                    {item.stockMax <= 5
                      ? `Осталось ${item.stockMax} шт.`
                      : `Доступно не более ${item.stockMax} шт.`}
                  </div>
                )}
              </div>
              <div className="cart-item-right">
                <div className="cart-item-price">{formatPrice(item.price * item.count)}</div>
                <button className="cart-item-remove" onClick={() => removeFromCart(item.id)}>✕</button>
              </div>
            </div>
          ))}
          <div className="cart-summary">
            <div className="cart-summary-row">
              <span>Итого</span>
              <span>{formatPrice(totalPrice)}</span>
            </div>
            <Link to="/checkout" className="cart-checkout">Перейти к оформлению</Link>
          </div>
        </div>
      )}
    </div>
  );
};

export default Cart;
