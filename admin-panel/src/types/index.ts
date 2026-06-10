// API Response wrapper
export interface ApiResponse<T = any> {
  status: string;
  data: T;
}

// Auth
export interface Admin {
  id: string;
  username: string;
  is_active: boolean;
  created_at: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  admin: Admin;
}

// Products
export interface Product {
  id: string;
  category_id: string;
  firm_id: string | null;
  name: string;
  description: string | null;
  price_kopeks: number;
  stock_quantity: number;
  is_active: boolean;
  created_at: string;
  images?: ProductImage[];
  tags?: Tag[];
}

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
  color: string;
  created_at: string;
}

// Categories
export interface Category {
  id: string;
  name: string;
  parent_id: string | null;
  image_url: string | null;
  is_active: boolean;
  created_at: string;
}

// Firms
export interface Firm {
  id: string;
  name: string;
  created_at: string;
}

// Orders
export interface Order {
  id: string;
  user_id: string;
  delivery_method_id: string | null;
  payment_method_id: string | null;
  promocode_id: string | null;
  status: 'created' | 'paid' | 'delivering' | 'completed' | 'cancelled';
  total_amount_kopeks: number;
  address_snapshot: string | null;
  created_at: string;
  items?: OrderItem[];
}

export interface OrderItem {
  id: string;
  order_id: string;
  product_id: string;
  product_name?: string;
  quantity: number;
  price_at_purchase_kopeks: number;
}

export interface OrderStatusHistory {
  id: string;
  order_id: string;
  status: string;
  comment: string | null;
  changed_at: string;
}

// Users
export interface User {
  id: string;
  telegram_id: number;
  username: string | null;
  first_name: string | null;
  last_name: string | null;
  phone: string | null;
  delivery_address: string | null;
  is_active: boolean;
  created_at: string;
}

// Promocodes
export interface Promocode {
  id: string;
  code: string;
  discount_percent: number;
  max_uses: number | null;
  current_uses: number;
  is_active: boolean;
  expires_at: string | null;
  created_at: string;
}

// Delivery Methods
export type DeliveryMethodKind = 'courier' | 'pickup' | 'pvz' | 'postal';

export interface DeliveryMethod {
  id: string;
  name: string;
  kind: DeliveryMethodKind;
  price_kopeks: number;
  description: string | null;
  is_active: boolean;
  created_at: string;
}

export const DELIVERY_KIND_LABELS: Record<DeliveryMethodKind, string> = {
  courier: 'Курьер (адрес доставки)',
  pickup: 'Самовывоз (без адреса)',
  pvz: 'ПВЗ / СДЭК',
  postal: 'Почта',
};

// Payment Methods
export interface PaymentMethod {
  id: string;
  name: string;
  description: string | null;
  is_active: boolean;
  created_at: string;
}

// Reviews
export interface Review {
  id: string;
  user_id: string;
  product_id: string;
  rating: number;
  comment: string | null;
  is_moderated: boolean;
  created_at: string;
}

// Notifications
export interface Notification {
  id: string;
  user_id: string;
  title: string | null;
  message: string | null;
  is_read: boolean;
  created_at: string;
}

// Paginated response
export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  limit: number;
  total_pages: number;
}
