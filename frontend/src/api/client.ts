import { formatApiError } from "../utils/apiError";
import {
  ApiResponse, Product, Category, Firm, PaginatedProducts,
  ProductRatingStat, Review, User, Order, CreateOrderInput, UpsertReviewInput,
  DeliveryMethod, PaymentMethod, Promocode, OrderStatusHistory,
  UserAddress, Notification as AppNotification,
  BasketItem as ServerBasketItem, FavoriteItem as ServerFavoriteItem, Tag,
} from "./types";

const BASE = "/api/v1";

function getInitData(): string {
  try {
    return (window as any).Telegram?.WebApp?.initData || "";
  } catch {
    return "";
  }
}

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE}${path}`, init);
  if (!res.ok) {
    const text = await res.text();
    const message = formatApiError(res.status, text);
    throw new Error(message);
  }
  const json: ApiResponse<T> = await res.json();
  return json.data ?? ([] as unknown as T);
}

async function authRequest<T>(path: string, init?: RequestInit): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(init?.headers as Record<string, string> || {}),
  };
  const initData = getInitData();
  if (initData) {
    headers["X-Telegram-Init-Data"] = initData;
  }
  return request<T>(path, { ...init, headers });
}

export const api = {
  getProducts: () => request<Product[]>("/products"),
  getProduct: (id: string) => request<Product>(`/products/${id}`),
  searchProducts: (params: Record<string, string>) => {
    const qs = new URLSearchParams(params).toString();
    return request<PaginatedProducts>(`/products/search?${qs}`);
  },
  getCategories: () => request<Category[]>("/categories"),
  getRootCategories: () => request<Category[]>("/categories/roots"),
  getSubcategories: (categoryId: string) => request<Category[]>(`/categories/${categoryId}/subcategories`),
  getCategory: (id: string) => request<Category>(`/categories/${id}`),
  getFirms: () => request<Firm[]>("/firms"),
  getTags: () => request<Tag[]>("/tags"),

  getProductRating: (productId: string) =>
    request<ProductRatingStat>(`/reviews/product/${productId}/stats`),
  getProductReviews: (productId: string) =>
    request<Review[]>(`/reviews/product/${productId}`),

  getDeliveryMethods: () => request<DeliveryMethod[]>("/delivery-methods"),
  getPaymentMethods: () => request<PaymentMethod[]>("/payment-methods"),
  validatePromocode: (code: string) =>
    request<Promocode>(`/promocodes/validate?code=${encodeURIComponent(code)}`),

  getUser: (telegramId: number) =>
    authRequest<User>(`/users/${telegramId}`),

  createUser: (data: { telegram_id: number; username?: string; first_name?: string; last_name?: string }) =>
    authRequest<User>("/users", { method: "POST", body: JSON.stringify(data) }),

  updateUserProfile: (telegramId: number, data: Partial<Pick<User, "username" | "first_name" | "last_name" | "phone" | "delivery_address">>) =>
    authRequest<null>(`/users/${telegramId}`, { method: "PUT", body: JSON.stringify(data) }),

  getAddresses: (telegramId: number) =>
    authRequest<UserAddress[]>(`/addresses/${telegramId}`),
  createAddress: (data: { user_id: number; label: string; full_address: string; is_default: boolean }) =>
    authRequest<UserAddress>("/addresses", { method: "POST", body: JSON.stringify(data) }),
  deleteAddress: (id: string) =>
    authRequest<null>(`/addresses/${id}`, { method: "DELETE" }),

  getBasket: (telegramId: number) =>
    authRequest<ServerBasketItem[]>(`/basket/${telegramId}`),
  addToBasket: (data: { user_id: number; product_id: string; quantity: number }) =>
    authRequest<null>("/basket", { method: "POST", body: JSON.stringify(data) }),
  updateBasketItem: (data: { user_id: number; product_id: string; quantity: number }) =>
    authRequest<null>("/basket", { method: "PUT", body: JSON.stringify(data) }),
  removeFromBasket: (telegramId: number, productId: string) =>
    authRequest<null>(`/basket/${telegramId}/${productId}`, { method: "DELETE" }),

  getFavorites: (telegramId: number) =>
    authRequest<ServerFavoriteItem[]>(`/favorites/${telegramId}`),
  addToFavorites: (data: { user_id: number; product_id: string }) =>
    authRequest<null>("/favorites", { method: "POST", body: JSON.stringify(data) }),
  removeFromFavorites: (telegramId: number, productId: string) =>
    authRequest<null>(`/favorites/${telegramId}/${productId}`, { method: "DELETE" }),

  createOrder: (data: CreateOrderInput) =>
    authRequest<Order>("/orders", { method: "POST", body: JSON.stringify(data) }),
  getUserOrders: (telegramId: number) =>
    authRequest<Order[]>(`/orders/user/${telegramId}`),
  getOrder: (id: string) =>
    authRequest<Order>(`/orders/${id}`),
  getOrderStatusHistory: (orderId: string) =>
    authRequest<OrderStatusHistory[]>(`/orders/${orderId}/history`),

  upsertReview: (data: UpsertReviewInput) =>
    authRequest<Review>("/reviews", { method: "POST", body: JSON.stringify(data) }),

  getNotifications: (telegramId: number) =>
    authRequest<AppNotification[]>(`/notifications/${telegramId}`),
  markNotificationRead: (id: string) =>
    authRequest<null>(`/notifications/${id}/read`, { method: "PATCH" }),
  markAllNotificationsRead: (telegramId: number) =>
    authRequest<null>(`/notifications/${telegramId}/read-all`, { method: "PATCH" }),
};
