import React, { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useCart } from "../context/CartContext";
import { useTelegram } from "../context/TelegramContext";
import { useToast } from "../context/ToastContext";
import { api } from "../api/client";
import { DeliveryMethod, PaymentMethod, UserAddress, formatPrice } from "../api/types";
import { formatFullAddress, parseFullAddress } from "../utils/addressParse";
import {
  buildAddressSnapshot,
  cityFieldLabel,
  cityFieldPlaceholder,
  needsStreetAddress,
  resolveDeliveryKind,
} from "./checkoutDeliveryUi";
import {
  FieldErrors,
  filterCityInput,
  filterPersonNameInput,
  firstErrorMessage,
  validateCheckoutForm,
} from "../utils/validation";
import { messageFromRequestError } from "../utils/apiError";
import "../style/checkout.css";

const Checkout: React.FC = () => {
  const { cartItems, totalPrice, clearCart } = useCart();
  const { user: tgUser } = useTelegram();
  const { showToast } = useToast();
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [phone, setPhone] = useState("");
  const [city, setCity] = useState("");
  const [streetAddress, setStreetAddress] = useState("");
  const [comment, setComment] = useState("");
  const [pickupPointLine, setPickupPointLine] = useState("");
  const [promoCode, setPromoCode] = useState("");
  const [promoDiscount, setPromoDiscount] = useState(0);
  const [promoApplied, setPromoApplied] = useState(false);
  const [submitting, setSubmitting] = useState(false);
  const [errors, setErrors] = useState<FieldErrors>({});

  const [deliveryMethods, setDeliveryMethods] = useState<DeliveryMethod[]>([]);
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [selectedDelivery, setSelectedDelivery] = useState("");
  const [selectedPayment, setSelectedPayment] = useState("");
  const [savedAddresses, setSavedAddresses] = useState<UserAddress[]>([]);
  const [selectedAddressId, setSelectedAddressId] = useState<string | null>(null);
  const [addressPrefilled, setAddressPrefilled] = useState(false);

  useEffect(() => {
    api.getDeliveryMethods().then(data => {
      const active = (data || []).filter(d => d.is_active);
      setDeliveryMethods(active);
      if (active.length > 0) setSelectedDelivery(active[0].id);
    }).catch(() => {});

    api.getPaymentMethods().then(data => {
      const active = (data || []).filter(p => p.is_active);
      setPaymentMethods(active);
      if (active.length > 0) setSelectedPayment(active[0].id);
    }).catch(() => {});
  }, []);

  useEffect(() => {
    if (!tgUser) return;
    Promise.all([
      api.getUser(tgUser.id).catch(() => null),
      api.getAddresses(tgUser.id).catch(() => [] as UserAddress[]),
    ]).then(([user, addrs]) => {
      if (user) {
        setName(`${user.first_name || ""} ${user.last_name || ""}`.trim());
        setPhone(user.phone || "");
      }
      const list = Array.isArray(addrs) ? addrs : [];
      if (list.length === 0 && user?.delivery_address?.trim()) {
        setSavedAddresses([{
          id: "__legacy__",
          user_id: user.id,
          full_address: user.delivery_address.trim(),
          is_default: true,
          created_at: "",
        }]);
      } else {
        setSavedAddresses(list);
      }
    });
  }, [tgUser]);

  const currentDelivery = deliveryMethods.find(d => d.id === selectedDelivery);
  const currentPayment = paymentMethods.find(p => p.id === selectedPayment);
  const deliveryKind = resolveDeliveryKind(currentDelivery);
  const courierLike = needsStreetAddress(deliveryKind);

  const clearFieldError = (field: string) => {
    setErrors(prev => {
      if (!prev[field]) return prev;
      const next = { ...prev };
      delete next[field];
      return next;
    });
  };

  const applySavedAddress = (addr: UserAddress) => {
    const { city: parsedCity, street } = parseFullAddress(addr.full_address);
    setCity(parsedCity);
    setStreetAddress(street);
    setSelectedAddressId(addr.id);
    clearFieldError("city");
    clearFieldError("streetAddress");
  };

  useEffect(() => {
    if (!courierLike || addressPrefilled || savedAddresses.length === 0) return;
    const def = savedAddresses.find(a => a.is_default) ?? savedAddresses[0];
    applySavedAddress(def);
    setAddressPrefilled(true);
  }, [courierLike, savedAddresses, addressPrefilled]);

  const deliveryCost = currentDelivery?.price_kopeks ?? 0;
  const subtotal = totalPrice + deliveryCost;
  const discountAmount = promoApplied ? Math.round(subtotal * promoDiscount / 100) : 0;
  const grandTotal = subtotal - discountAmount;

  const handleApplyPromo = async () => {
    if (!promoCode.trim()) return;
    try {
      const promo = await api.validatePromocode(promoCode.trim());
      setPromoDiscount(promo.discount_percent);
      setPromoApplied(true);
      showToast(`Промокод применён: −${promo.discount_percent}%`);
    } catch {
      showToast("Промокод недействителен", "error");
      setPromoApplied(false);
      setPromoDiscount(0);
    }
  };

  const handleSelectDelivery = (opt: DeliveryMethod) => {
    setSelectedDelivery(opt.id);
    setPickupPointLine("");
    const kind = resolveDeliveryKind(opt);
    if (!needsStreetAddress(kind)) {
      setStreetAddress("");
      setSelectedAddressId(null);
    } else if (savedAddresses.length > 0 && !streetAddress && !city) {
      const def = savedAddresses.find(a => a.is_default) ?? savedAddresses[0];
      applySavedAddress(def);
    }
    clearFieldError("delivery");
    clearFieldError("streetAddress");
    clearFieldError("pickupPointLine");
  };

  const handleOrder = async () => {
    if (!tgUser) {
      showToast("Откройте через Telegram", "error");
      return;
    }

    const validationErrors = validateCheckoutForm({
      name,
      phone,
      city,
      streetAddress,
      pickupPointLine,
      courierLike,
      deliveryKind,
      selectedDelivery,
      selectedPayment,
      cartEmpty: cartItems.length === 0,
    });
    setErrors(validationErrors);
    const msg = firstErrorMessage(validationErrors);
    if (msg) {
      showToast(msg, "error");
      return;
    }

    const snapshot = buildAddressSnapshot({
      kind: deliveryKind,
      methodLabel: currentDelivery?.name ?? "",
      city,
      streetAddress,
      comment,
      pickupPointLine,
    });

    setSubmitting(true);
    try {
      for (const item of cartItems) {
        const product = await api.getProduct(item.id);
        if (item.count > product.stock_quantity) {
          showToast(
            `«${item.title}»: доступно только ${product.stock_quantity} шт.`,
            "error"
          );
          setSubmitting(false);
          return;
        }
      }

      await api.updateUserProfile(tgUser.id, {
        first_name: name.split(" ")[0] || undefined,
        last_name: name.split(" ").slice(1).join(" ") || undefined,
        phone: phone || undefined,
        delivery_address: courierLike
          ? formatFullAddress(city, streetAddress) || undefined
          : undefined,
      });

      await api.createOrder({
        user_id: tgUser.id,
        delivery_method_id: selectedDelivery,
        payment_method_id: selectedPayment,
        promo_code: promoApplied ? promoCode.trim() : undefined,
        delivery_address: snapshot,
        items: cartItems.map(item => ({
          product_id: item.id,
          quantity: item.count,
        })),
      });

      clearCart();
      showToast("Заказ оформлен!");
      navigate("/orders");
    } catch (err: unknown) {
      showToast(messageFromRequestError(err), "error");
    }
    setSubmitting(false);
  };

  if (cartItems.length === 0) {
    return (
      <div className="checkout-container">
        <h1 className="checkout-title">Оформление заказа</h1>
        <div className="checkout-empty">
          <p>Корзина пуста</p>
          <button className="checkout-back-btn" onClick={() => navigate("/")}>На главную</button>
        </div>
      </div>
    );
  }

  return (
    <div className="checkout-container">
      <h1 className="checkout-title">Оформление заказа</h1>

      <section className="checkout-section">
        <h2 className="checkout-subtitle">Покупатель</h2>
        <label className={errors.name ? "checkout-label-invalid" : ""}>
          Ф.И.О.*
          <input
            value={name}
            onChange={e => { setName(filterPersonNameInput(e.target.value)); clearFieldError("name"); }}
            className={errors.name ? "input-invalid" : ""}
          />
          {errors.name && <span className="field-error">{errors.name}</span>}
        </label>
        <label className={errors.phone ? "checkout-label-invalid" : ""}>
          Мобильный телефон*
          <input
            value={phone}
            onChange={e => { setPhone(e.target.value); clearFieldError("phone"); }}
            type="tel"
            placeholder="+375 29 123-45-67"
            className={errors.phone ? "input-invalid" : ""}
          />
          {errors.phone && <span className="field-error">{errors.phone}</span>}
        </label>
      </section>

      <section className="checkout-section">
        <h2 className="checkout-subtitle">Способ доставки</h2>
        {errors.delivery && <span className="field-error">{errors.delivery}</span>}
        <div className="delivery-options">
          {deliveryMethods.map(opt => (
            <label
              key={opt.id}
              className={`delivery-option${selectedDelivery === opt.id ? " selected" : ""}`}
            >
              <div className="delivery-text">
                <div className="delivery-label">{opt.name}</div>
                <div className="delivery-price">
                  {opt.price_kopeks === 0
                    ? formatPrice(0)
                    : formatPrice(opt.price_kopeks)}
                  {opt.description && opt.price_kopeks === 0 && (
                    <span className="delivery-price-hint"> ({opt.description})</span>
                  )}
                </div>
                {opt.description && opt.price_kopeks !== 0 && (
                  <div className="delivery-note">{opt.description}</div>
                )}
              </div>
              <input
                type="radio"
                name="delivery"
                value={opt.id}
                checked={selectedDelivery === opt.id}
                onChange={() => handleSelectDelivery(opt)}
              />
            </label>
          ))}
        </div>

        <label className={errors.city ? "checkout-label-invalid" : ""}>
          {cityFieldLabel(deliveryKind)}
          <input
            value={city}
            onChange={e => {
              setCity(filterCityInput(e.target.value));
              setSelectedAddressId(null);
              clearFieldError("city");
            }}
            placeholder={cityFieldPlaceholder(deliveryKind)}
            className={errors.city ? "input-invalid" : ""}
          />
          {errors.city && <span className="field-error">{errors.city}</span>}
        </label>

        {deliveryKind === "pvz" && (
          <div className="checkout-pvz-block">
            <p className="checkout-hint">
              Укажите пункт выдачи СДЭК — название и адрес (можно скопировать из приложения СДЭК).
            </p>
            <label className={errors.pickupPointLine ? "checkout-label-invalid" : ""}>
              Пункт выдачи (ПВЗ)*
              <textarea
                className={`checkout-textarea${errors.pickupPointLine ? " input-invalid" : ""}`}
                value={pickupPointLine}
                onChange={e => { setPickupPointLine(e.target.value); clearFieldError("pickupPointLine"); }}
                placeholder="СДЭК ПВЗ, ул. Примерная 1"
                rows={3}
              />
              {errors.pickupPointLine && <span className="field-error">{errors.pickupPointLine}</span>}
            </label>
          </div>
        )}

        {courierLike && savedAddresses.length > 0 && (
          <div className="checkout-saved-addresses">
            <span className="checkout-saved-label">Адрес из профиля</span>
            <div className="checkout-saved-list">
              {savedAddresses.map(addr => (
                <button
                  key={addr.id}
                  type="button"
                  className={`checkout-saved-item${selectedAddressId === addr.id ? " selected" : ""}`}
                  onClick={() => applySavedAddress(addr)}
                >
                  {addr.full_address}
                  {addr.is_default && <span className="checkout-saved-badge">основной</span>}
                </button>
              ))}
            </div>
          </div>
        )}

        {courierLike && (
          <>
            <label className={errors.streetAddress ? "checkout-label-invalid" : ""}>
              Улица, дом, подъезд*
              <textarea
                className={`checkout-textarea${errors.streetAddress ? " input-invalid" : ""}`}
                value={streetAddress}
                onChange={e => {
                  setStreetAddress(e.target.value);
                  setSelectedAddressId(null);
                  clearFieldError("streetAddress");
                }}
                placeholder="Мира 18, подъезд 2, кв. 45"
                rows={3}
              />
              {errors.streetAddress && <span className="field-error">{errors.streetAddress}</span>}
            </label>
            <p className="checkout-hint">
              Город указан выше — здесь только улица, дом и детали доставки.
            </p>
            <label>
              Комментарий к заказу
              <input value={comment} onChange={e => setComment(e.target.value)} placeholder="По желанию" />
            </label>
          </>
        )}

        {deliveryKind === "pickup" && (
          <p className="checkout-hint">
            Самовывоз: адрес магазина указан в описании способа доставки. Отдельный адрес не нужен.
          </p>
        )}

        {(deliveryKind === "pickup" || deliveryKind === "pvz") && (
          <label>
            Комментарий к заказу
            <input value={comment} onChange={e => setComment(e.target.value)} placeholder="По желанию" />
          </label>
        )}
      </section>

      <section className="checkout-section">
        <h2 className="checkout-subtitle">Способ оплаты</h2>
        <div className="delivery-options">
          {paymentMethods.map(opt => (
            <label
              key={opt.id}
              className={`delivery-option${selectedPayment === opt.id ? " selected" : ""}`}
            >
              <div className="delivery-text">
                <div className="delivery-label">{opt.name}</div>
                {opt.description && <div className="delivery-note">{opt.description}</div>}
              </div>
              <input
                type="radio"
                name="payment"
                value={opt.id}
                checked={selectedPayment === opt.id}
                onChange={() => { setSelectedPayment(opt.id); clearFieldError("delivery"); }}
              />
            </label>
          ))}
        </div>
      </section>

      <section className="checkout-section">
        <h2 className="checkout-subtitle">Промокод</h2>
        <div className="promo-row">
          <input
            value={promoCode}
            onChange={e => { setPromoCode(e.target.value); setPromoApplied(false); setPromoDiscount(0); }}
            placeholder="Введите промокод"
            className="promo-input"
          />
          <button className="promo-btn" onClick={handleApplyPromo} type="button">
            Применить
          </button>
        </div>
        {promoApplied && <div className="promo-applied">Скидка {promoDiscount}% применена</div>}
      </section>

      <section className="checkout-section">
        <h2 className="checkout-subtitle">Ваш заказ</h2>
        {cartItems.map(item => (
          <div className="checkout-item" key={item.id}>
            <div className="checkout-item-left">
              <span>{item.title} × {item.count}</span>
              {item.count > item.stockMax && (
                <span className="checkout-item-warning">
                  Доступно только {item.stockMax} шт.
                </span>
              )}
            </div>
            <span>{formatPrice(item.price * item.count)}</span>
          </div>
        ))}
      </section>

      <section className="checkout-summary-section">
        <div className="checkout-summary-block">
          <div className="checkout-summary-row">
            <span>Товары:</span>
            <span>{formatPrice(totalPrice)}</span>
          </div>
          <div className="checkout-summary-row">
            <span>Доставка:</span>
            <span>
              {deliveryCost === 0 ? "бесплатно" : formatPrice(deliveryCost)}
            </span>
          </div>
          {currentPayment && (
            <div className="checkout-summary-row checkout-summary-pay">
              <span>Оплата:</span>
              <span>{currentPayment.name}</span>
            </div>
          )}
          {promoApplied && (
            <div className="checkout-summary-row">
              <span>Скидка ({promoDiscount}%):</span>
              <span>−{formatPrice(discountAmount)}</span>
            </div>
          )}
          <div className="checkout-summary-row">
            <span>Итого:</span>
            <span className="checkout-summary-total">{formatPrice(grandTotal)}</span>
          </div>
          <button
            className="checkout-order-btn"
            onClick={handleOrder}
            disabled={submitting}
          >
            {submitting ? "Оформление..." : "Оформить заказ"}
          </button>
        </div>
        <div className="checkout-summary-note">
          Нажимая на кнопку, вы соглашаетесь на обработку персональных данных
        </div>
      </section>
    </div>
  );
};

export default Checkout;
