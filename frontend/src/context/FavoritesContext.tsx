import React, { createContext, useContext, useState, useEffect, useCallback, ReactNode } from "react";
import { useTelegram } from "./TelegramContext";
import { api } from "../api/client";

type FavoritesContextType = {
  favoriteIds: Set<string>;
  isFavorite: (productId: string) => boolean;
  toggleFavorite: (productId: string) => void;
  clearFavorites: () => void;
  count: number;
};

const FavoritesContext = createContext<FavoritesContextType | undefined>(undefined);

export const FavoritesProvider = ({ children }: { children: ReactNode }) => {
  const { user: tgUser } = useTelegram();
  const [favoriteIds, setFavoriteIds] = useState<Set<string>>(new Set());

  useEffect(() => {
    if (!tgUser) return;
    api.getFavorites(tgUser.id)
      .then(items => {
        const ids = new Set((items || []).map(i => i.product_id));
        setFavoriteIds(ids);
      })
      .catch(() => {});
  }, [tgUser]);

  const isFavorite = useCallback((productId: string) => favoriteIds.has(productId), [favoriteIds]);

  const toggleFavorite = useCallback((productId: string) => {
    setFavoriteIds(prev => {
      const next = new Set(prev);
      if (next.has(productId)) {
        next.delete(productId);
        if (tgUser) api.removeFromFavorites(tgUser.id, productId).catch(() => {});
      } else {
        next.add(productId);
        if (tgUser) api.addToFavorites({ user_id: tgUser.id, product_id: productId }).catch(() => {});
      }
      return next;
    });
  }, [tgUser]);

  const clearFavorites = useCallback(() => {
    favoriteIds.forEach(id => {
      if (tgUser) api.removeFromFavorites(tgUser.id, id).catch(() => {});
    });
    setFavoriteIds(new Set());
  }, [tgUser, favoriteIds]);

  return (
    <FavoritesContext.Provider value={{ favoriteIds, isFavorite, toggleFavorite, clearFavorites, count: favoriteIds.size }}>
      {children}
    </FavoritesContext.Provider>
  );
};

export const useFavorites = () => {
  const ctx = useContext(FavoritesContext);
  if (!ctx) throw new Error("useFavorites must be used within FavoritesProvider");
  return ctx;
};
