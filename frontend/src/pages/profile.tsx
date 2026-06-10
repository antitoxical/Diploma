import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { useTelegram } from "../context/TelegramContext";
import { useTheme } from "../context/ThemeContext";
import { useToast } from "../context/ToastContext";
import { api } from "../api/client";
import { UserAddress } from "../api/types";
import {
  FieldErrors,
  firstErrorMessage,
  filterPersonNameInput,
  validateNewAddress,
  validateProfileForm,
} from "../utils/validation";
import "../style/profile.css";

const Profile: React.FC = () => {
  const { user: tgUser } = useTelegram();
  const { mode, setMode } = useTheme();
  const { showToast } = useToast();
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);

  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [phone, setPhone] = useState("");
  const [errors, setErrors] = useState<FieldErrors>({});
  const [addressError, setAddressError] = useState<string | null>(null);

  const [addresses, setAddresses] = useState<UserAddress[]>([]);
  const [newAddress, setNewAddress] = useState("");

  useEffect(() => {
    if (!tgUser) { setLoading(false); return; }
    api.getUser(tgUser.id)
      .then(u => {
        setFirstName(u.first_name || tgUser.first_name || "");
        setLastName(u.last_name || tgUser.last_name || "");
        setPhone(u.phone || "");
      })
      .catch(() => {
        api.createUser({
          telegram_id: tgUser.id,
          username: tgUser.username,
          first_name: tgUser.first_name,
          last_name: tgUser.last_name,
        }).then(u => {
          setFirstName(u.first_name || tgUser.first_name || "");
          setLastName(u.last_name || "");
        }).catch(() => {});
      })
      .finally(() => setLoading(false));

    api.getAddresses(tgUser.id)
      .then(data => setAddresses(Array.isArray(data) ? data : []))
      .catch(() => {});
  }, [tgUser]);

  const clearFieldError = (field: string) => {
    setErrors(prev => {
      if (!prev[field]) return prev;
      const next = { ...prev };
      delete next[field];
      return next;
    });
  };

  const handleSave = async () => {
    if (!tgUser) return;
    const validationErrors = validateProfileForm({ firstName, lastName, phone });
    setErrors(validationErrors);
    const msg = firstErrorMessage(validationErrors);
    if (msg) {
      showToast(msg, "error");
      return;
    }

    setSaving(true);
    try {
      await api.updateUserProfile(tgUser.id, {
        first_name: firstName.trim(),
        last_name: lastName.trim() || undefined,
        phone: phone.trim(),
      });
      showToast("Профиль сохранён");
    } catch {
      showToast("Ошибка сохранения", "error");
    }
    setSaving(false);
  };

  const handleAddAddress = async () => {
    if (!tgUser) return;
    const err = validateNewAddress(newAddress, addresses.length);
    setAddressError(err);
    if (err) {
      showToast(err, "error");
      return;
    }

    try {
      const addr = await api.createAddress({
        user_id: tgUser.id,
        label: "",
        full_address: newAddress.trim(),
        is_default: addresses.length === 0,
      });
      setAddresses(prev => [...prev, addr]);
      setNewAddress("");
      setAddressError(null);
      showToast("Адрес добавлен");
    } catch {
      showToast("Ошибка добавления адреса", "error");
    }
  };

  const handleDeleteAddress = async (id: string) => {
    try {
      await api.deleteAddress(id);
      setAddresses(prev => prev.filter(a => a.id !== id));
      showToast("Адрес удалён");
    } catch {
      showToast("Ошибка удаления", "error");
    }
  };

  if (loading) {
    return (
      <div className="profile-container">
        <p className="profile-loading">Загрузка...</p>
      </div>
    );
  }

  if (!tgUser) {
    return (
      <div className="profile-container">
        <h1 className="profile-title">Профиль</h1>
        <p className="profile-hint">Откройте приложение через Telegram для доступа к профилю</p>
        <div className="profile-theme-section">
          <span className="profile-theme-label">Тема</span>
          <div className="theme-toggle">
            {(["light", "system", "dark"] as const).map(m => (
              <button
                key={m}
                className={`theme-toggle-btn ${mode === m ? "active" : ""}`}
                onClick={() => setMode(m)}
              >
                {m === "light" ? "Светлая" : m === "dark" ? "Тёмная" : "Авто"}
              </button>
            ))}
          </div>
        </div>
      </div>
    );
  }

  const avatarUrl = tgUser.photo_url;
  const avatarLetter = (firstName || tgUser.first_name || "U").charAt(0).toUpperCase();

  return (
    <div className="profile-container">
      <h1 className="profile-title">Профиль</h1>

      {avatarUrl ? (
        <img src={avatarUrl} alt="avatar" className="profile-avatar-img" />
      ) : (
        <div className="profile-avatar">{avatarLetter}</div>
      )}
      <div className="profile-username">@{tgUser.username || "user"}</div>

      <div className="profile-form">
        <label className={`profile-label${errors.firstName ? " profile-label-invalid" : ""}`}>
          Имя
          <input
            value={firstName}
            onChange={e => { setFirstName(filterPersonNameInput(e.target.value)); clearFieldError("firstName"); }}
            placeholder="Имя"
            className={errors.firstName ? "input-invalid" : ""}
          />
          {errors.firstName && <span className="field-error">{errors.firstName}</span>}
        </label>
        <label className={`profile-label${errors.lastName ? " profile-label-invalid" : ""}`}>
          Фамилия
          <input
            value={lastName}
            onChange={e => { setLastName(filterPersonNameInput(e.target.value)); clearFieldError("lastName"); }}
            placeholder="Фамилия"
            className={errors.lastName ? "input-invalid" : ""}
          />
          {errors.lastName && <span className="field-error">{errors.lastName}</span>}
        </label>
        <label className={`profile-label${errors.phone ? " profile-label-invalid" : ""}`}>
          Телефон
          <input
            value={phone}
            onChange={e => { setPhone(e.target.value); clearFieldError("phone"); }}
            placeholder="+375 29 123-45-67"
            type="tel"
            className={errors.phone ? "input-invalid" : ""}
          />
          {errors.phone && <span className="field-error">{errors.phone}</span>}
        </label>
        <button className="profile-save" onClick={handleSave} disabled={saving}>
          {saving ? "Сохранение..." : "Сохранить"}
        </button>
      </div>

      <div className="profile-addresses">
        <h3>Мои адреса</h3>
        <p className="profile-address-hint">Формат «Город, улица» — при курьерской доставке подставится в оформление заказа.</p>
        {addresses.map(a => (
          <div key={a.id} className="profile-address-item">
            <span>{a.full_address}</span>
            {a.is_default && <span className="profile-address-default">основной</span>}
            <button onClick={() => handleDeleteAddress(a.id)} className="profile-address-delete">✕</button>
          </div>
        ))}
        <div className="profile-address-add">
          <div className="profile-address-input-wrap">
            <input
              value={newAddress}
              onChange={e => { setNewAddress(e.target.value); setAddressError(null); }}
              placeholder="Город, улица, дом (например: Смоленск, Мира 18)"
              className={addressError ? "input-invalid" : ""}
            />
            {addressError && <span className="field-error">{addressError}</span>}
          </div>
          <button onClick={handleAddAddress}>+</button>
        </div>
      </div>

      <div className="profile-theme-section">
        <span className="profile-theme-label">Тема оформления</span>
        <div className="theme-toggle">
          <button
            className={`theme-toggle-btn ${mode === "light" ? "active" : ""}`}
            onClick={() => setMode("light")}
          >
            Светлая
          </button>
          <button
            className={`theme-toggle-btn ${mode === "system" ? "active" : ""}`}
            onClick={() => setMode("system")}
          >
            Авто
          </button>
          <button
            className={`theme-toggle-btn ${mode === "dark" ? "active" : ""}`}
            onClick={() => setMode("dark")}
          >
            Тёмная
          </button>
        </div>
      </div>

      <div className="profile-links">
        <Link to="/orders" className="profile-link-item">
          <span>Мои заказы</span>
          <span className="profile-link-arrow">&gt;</span>
        </Link>
        <Link to="/favorites" className="profile-link-item">
          <span>Избранное</span>
          <span className="profile-link-arrow">&gt;</span>
        </Link>
        <Link to="/notifications" className="profile-link-item">
          <span>Уведомления</span>
          <span className="profile-link-arrow">&gt;</span>
        </Link>
      </div>
    </div>
  );
};

export default Profile;
