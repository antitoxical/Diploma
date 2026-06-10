import { kopeksToRubles, rublesToKopeks } from './api';
import type { Product, ProductImage } from '../types';

/** Поля формы редактирования товара (не совпадают с ответом API). */
export type ProductFormValues = {
  name: string;
  category_id: string;
  firm_id?: string | null;
  description?: string | null;
  price_rubles: number;
  stock_quantity: number;
  is_active: boolean;
  images?: string[];
  tag_ids?: string[];
};

export function mapProductToFormValues(product: Product): ProductFormValues {
  return {
    name: product.name,
    description: product.description,
    category_id: product.category_id,
    firm_id: product.firm_id,
    price_rubles: kopeksToRubles(product.price_kopeks),
    stock_quantity: product.stock_quantity,
    is_active: product.is_active,
    tag_ids: product.tags?.map((t) => t.id) ?? [],
    images: product.images?.map((img) => img.image_url) ?? [],
  };
}

export function toImageUrlList(images: unknown, fallback?: ProductFormValues): string[] {
  if (Array.isArray(images) && images.length > 0) {
    return images
      .map((item) => {
        if (typeof item === 'string') return item;
        if (item && typeof item === 'object' && 'image_url' in item) {
          return String((item as ProductImage).image_url);
        }
        return '';
      })
      .filter(Boolean);
  }
  return fallback?.images ?? [];
}

export function productFormValuesToUpdatePayload(
  values: ProductFormValues,
  fallback?: ProductFormValues,
) {
  const priceRubles = values.price_rubles ?? fallback?.price_rubles ?? 0;

  return {
    name: values.name,
    category_id: values.category_id,
    firm_id: values.firm_id ?? null,
    description: values.description ?? null,
    price_kopeks: rublesToKopeks(priceRubles),
    stock_quantity: values.stock_quantity,
    is_active: values.is_active,
    image_urls: toImageUrlList(values.images, fallback),
  };
}
