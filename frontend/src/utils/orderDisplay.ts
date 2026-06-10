import { OrderItemRow } from "../api/types";

export function orderItemTitle(item: OrderItemRow): string {
  if (item.product_name?.trim()) return item.product_name.trim();
  return `Товар #${item.product_id.slice(0, 8)}`;
}

export function orderTotalQuantity(products: OrderItemRow[] = []): number {
  return products.reduce((sum, item) => sum + item.quantity, 0);
}

export function orderProductsPreview(products: OrderItemRow[] = []): string {
  if (products.length === 0) return "Нет товаров";
  const first = orderItemTitle(products[0]);
  if (products.length === 1) {
    return `${first} × ${products[0].quantity}`;
  }
  const rest = products.length - 1;
  return `${first} и ещё ${rest} поз.`;
}

export function formatOrderAddress(snapshot?: string): string[] {
  if (!snapshot?.trim()) return [];
  return snapshot
    .split("\n")
    .map(line => line.trim())
    .filter(Boolean);
}
