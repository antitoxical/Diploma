import React, { createContext, useContext, useState, useEffect, useCallback, useRef, ReactNode } from "react";
import { useTelegram } from "./TelegramContext";
import { api } from "../api/client";
import { getMainImage } from "../api/types";

export type CartItem = {
  id: string;
  title: string;
  image: string;
  price: number;
  count: number;
  stockMax: number;
};

type CartContextType = {
  cartItems: CartItem[];
  addToCart: (item: Omit<CartItem, "count">) => boolean;
  removeFromCart: (id: string) => void;
  increment: (id: string) => boolean;
  decrement: (id: string) => void;
  clearCart: () => void;
  refreshStock: () => Promise<void>;
  totalPrice: number;
  loading: boolean;
};

const CartContext = createContext<CartContextType | undefined>(undefined);

async function buildCartItem(
  si: { product_id: string; quantity: number },
  tgUserId: number
): Promise<CartItem | null> {
  try {
    const product = await api.getProduct(si.product_id);
    const stockMax = Math.max(0, product.stock_quantity);
    if (stockMax <= 0) {
      api.removeFromBasket(tgUserId, si.product_id).catch(() => {});
      return null;
    }

    const count = Math.min(si.quantity, stockMax);
    if (count !== si.quantity) {
      api.updateBasketItem({
        user_id: tgUserId,
        product_id: si.product_id,
        quantity: count,
      }).catch(() => {});
    }

    return {
      id: si.product_id,
      title: product.name,
      image: getMainImage(product),
      price: product.price_kopeks,
      count,
      stockMax,
    };
  } catch {
    return null;
  }
}

export const CartProvider = ({ children }: { children: ReactNode }) => {
  const { user: tgUser } = useTelegram();
  const [cartItems, setCartItems] = useState<CartItem[]>([]);
  const [loading, setLoading] = useState(false);
  const cartItemsRef = useRef(cartItems);
  cartItemsRef.current = cartItems;

  const loadCart = useCallback(async () => {
    if (!tgUser) {
      setCartItems([]);
      return;
    }
    setLoading(true);
    try {
      const serverItems = await api.getBasket(tgUser.id);
      if (!serverItems?.length) {
        setCartItems([]);
        return;
      }
      const items = await Promise.all(
        serverItems.map(si => buildCartItem(si, tgUser.id))
      );
      setCartItems(items.filter((i): i is CartItem => i !== null));
    } catch {
      setCartItems([]);
    } finally {
      setLoading(false);
    }
  }, [tgUser]);

  useEffect(() => {
    loadCart();
  }, [loadCart]);

  const refreshStock = useCallback(async () => {
    if (!tgUser) return;
    const current = cartItemsRef.current;
    if (current.length === 0) return;

    const updated = await Promise.all(
      current.map(async (item) => {
        try {
          const product = await api.getProduct(item.id);
          const stockMax = Math.max(0, product.stock_quantity);
          if (stockMax <= 0) {
            api.removeFromBasket(tgUser.id, item.id).catch(() => {});
            return null;
          }
          const count = Math.min(item.count, stockMax);
          if (count !== item.count) {
            api.updateBasketItem({
              user_id: tgUser.id,
              product_id: item.id,
              quantity: count,
            }).catch(() => {});
          }
          return {
            ...item,
            title: product.name,
            price: product.price_kopeks,
            count,
            stockMax,
          };
        } catch {
          return item;
        }
      })
    );

    setCartItems(updated.filter((i): i is CartItem => i !== null));
  }, [tgUser]);

  const syncAdd = useCallback((productId: string, quantity: number) => {
    if (tgUser) {
      api.addToBasket({ user_id: tgUser.id, product_id: productId, quantity }).catch(() => {});
    }
  }, [tgUser]);

  const syncUpdate = useCallback((productId: string, quantity: number) => {
    if (tgUser) {
      api.updateBasketItem({ user_id: tgUser.id, product_id: productId, quantity }).catch(() => {});
    }
  }, [tgUser]);

  const syncRemove = useCallback((productId: string) => {
    if (tgUser) {
      api.removeFromBasket(tgUser.id, productId).catch(() => {});
    }
  }, [tgUser]);

  const addToCart = (item: Omit<CartItem, "count">): boolean => {
    if (item.stockMax <= 0) return false;

    let added = false;
    setCartItems(prev => {
      const existing = prev.find(i => i.id === item.id);
      const stockMax = existing?.stockMax ?? item.stockMax;

      if (existing) {
        if (existing.count >= stockMax) return prev;
        const newCount = existing.count + 1;
        syncUpdate(item.id, newCount);
        added = true;
        return prev.map(i =>
          i.id === item.id ? { ...i, count: newCount, stockMax } : i
        );
      }

      syncAdd(item.id, 1);
      added = true;
      return [...prev, { ...item, count: 1, stockMax }];
    });
    return added;
  };

  const removeFromCart = (id: string) => {
    syncRemove(id);
    setCartItems(prev => prev.filter(i => i.id !== id));
  };

  const increment = (id: string): boolean => {
    let ok = false;
    setCartItems(prev => {
      const item = prev.find(i => i.id === id);
      if (!item || item.count >= item.stockMax) return prev;
      syncUpdate(id, item.count + 1);
      ok = true;
      return prev.map(i => (i.id === id ? { ...i, count: i.count + 1 } : i));
    });
    return ok;
  };

  const decrement = (id: string) => {
    setCartItems(prev => {
      const item = prev.find(i => i.id === id);
      if (item && item.count <= 1) {
        syncRemove(id);
        return prev.filter(i => i.id !== id);
      }
      if (item) syncUpdate(id, item.count - 1);
      return prev.map(i => (i.id === id ? { ...i, count: i.count - 1 } : i));
    });
  };

  const clearCart = () => {
    cartItems.forEach(i => syncRemove(i.id));
    setCartItems([]);
  };

  const totalPrice = cartItems.reduce((sum, i) => sum + i.price * i.count, 0);

  return (
    <CartContext.Provider value={{
      cartItems,
      addToCart,
      removeFromCart,
      increment,
      decrement,
      clearCart,
      refreshStock,
      totalPrice,
      loading,
    }}>
      {children}
    </CartContext.Provider>
  );
};

export const useCart = () => {
  const ctx = useContext(CartContext);
  if (!ctx) throw new Error("useCart must be used within CartProvider");
  return ctx;
};
