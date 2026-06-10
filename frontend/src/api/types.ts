export interface ProductImage {
  id: string;
  product_id: string;
  image_url: string;
  is_main: boolean;
  sort_order: number;
}

export interface Tag {
  id: string;
  name: string;
  color?: string;
  created_at: string;
}

export interface Product {
  id: string;
  category_id: string;
  firm_id: string;
  name: string;
  description: string;
  price_kopeks: number;
  stock_quantity: number;
  is_active: boolean;
  created_at: string;
  images: ProductImage[];
  tags: Tag[];
}

export interface Category {
  id: string;
  name: string;
  parent_id?: string;
  image_url: string;
  is_active: boolean;
  created_at: string;
}

export interface Firm {
  id: string;
  name: string;
  created_at: string;
}

export interface PaginatedProducts {
  products: Product[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}

export interface ProductRatingStat {
  product_id: string;
  avg_rating: number;
  review_count: number;
}

export interface User {
  id: string;
  telegram_id: number;
  username?: string;
  first_name?: string;
  last_name?: string;
  phone?: string;
  delivery_address?: string;
  is_active: boolean;
  created_at: string;
}

export interface UserAddress {
  id: string;
  user_id: string;
  label?: string;
  full_address: string;
  is_default: boolean;
  created_at: string;
}

export interface Review {
  id: string;
  user_id: string;
  product_id: string;
  rating: number;
  comment: string;
  is_moderated: boolean;
  created_at: string;
}

export interface UpsertReviewInput {
  user_id: number;
  product_id: string;
  rating: number;
  comment?: string;
}

export type DeliveryMethodKind = "courier" | "pickup" | "pvz" | "postal";

export interface DeliveryMethod {
  id: string;
  name: string;
  kind: DeliveryMethodKind;
  price_kopeks: number;
  description?: string;
  is_active: boolean;
  created_at: string;
}

export interface PaymentMethod {
  id: string;
  name: string;
  description?: string;
  is_active: boolean;
  created_at: string;
}

export interface Promocode {
  id: string;
  code: string;
  discount_percent: number;
  max_uses?: number;
  current_uses: number;
  is_active: boolean;
  expires_at?: string;
  created_at: string;
}

export interface OrderItemRow {
  id: string;
  order_id: string;
  product_id: string;
  product_name?: string;
  quantity: number;
  price_at_purchase_kopeks: number;
}

export interface Order {
  id: string;
  user_id: string;
  telegram_id: number;
  delivery_method_id?: string;
  payment_method_id?: string;
  delivery_method_name?: string;
  payment_method_name?: string;
  promocode_id?: string;
  status: string;
  total_amount_kopeks: number;
  address_snapshot?: string;
  created_at: string;
  products: OrderItemRow[];
}

export interface OrderStatusHistory {
  id: string;
  order_id: string;
  status: string;
  comment?: string;
  changed_at: string;
}

export interface CreateOrderInput {
  user_id: number;
  delivery_method_id: string;
  payment_method_id: string;
  promo_code?: string;
  delivery_address: string;
  items: { product_id: string; quantity: number }[];
}

export interface Notification {
  id: string;
  user_id: string;
  title?: string;
  message?: string;
  is_read: boolean;
  created_at: string;
}

export interface BasketItem {
  id: string;
  user_id: string;
  product_id: string;
  quantity: number;
  created_at: string;
}

export interface FavoriteItem {
  id: string;
  user_id: string;
  product_id: string;
  created_at: string;
}

export interface ApiResponse<T> {
  status: string;
  data: T;
}

export const PLACEHOLDER_IMAGE = "/images/placeholder.svg";

export function rewriteExternalImageUrl(url: string): string {
  if (!url || url.startsWith("/")) return url;
  try {
    const u = new URL(url);
    const host = u.hostname.toLowerCase();
    if (host === "placehold.co" || host === "www.placehold.co") {
      return `/ext-img${u.pathname}${u.search}`;
    }
  } catch {
  }
  return url;
}

export function getMainImage(product: Product): string {
  if (!product.images || product.images.length === 0) return PLACEHOLDER_IMAGE;
  const main = product.images.find(i => i.is_main);
  const raw = (main || product.images[0]).image_url;
  return rewriteExternalImageUrl(raw);
}

export function getAllImages(product: Product): string[] {
  if (!product.images || product.images.length === 0) return [PLACEHOLDER_IMAGE];
  return product.images
    .sort((a, b) => a.sort_order - b.sort_order)
    .map(i => rewriteExternalImageUrl(i.image_url));
}

export function kopeksToRubles(kopeks: number): number {
  return kopeks / 100;
}

export function formatPrice(kopeks: number): string {
  const major = kopeksToRubles(kopeks);
  return `${major.toLocaleString("be-BY", { minimumFractionDigits: 2, maximumFractionDigits: 2 })} Br`;
}

export function orderStatusLabel(status: string): string {
  const map: Record<string, string> = {
    created: "Создан",
    paid: "Оплачен",
    delivering: "Доставляется",
    completed: "Завершён",
    cancelled: "Отменён",
  };
  return map[status] || status;
}

export function orderStatusColor(status: string): string {
  const map: Record<string, string> = {
    created: "#ff9800",
    paid: "#2196f3",
    delivering: "#9c27b0",
    completed: "#4caf50",
    cancelled: "#e53935",
  };
  return map[status] || "#888";
}
