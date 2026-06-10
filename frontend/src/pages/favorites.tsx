import React, { useEffect, useState } from "react";
import { useFavorites } from "../context/FavoritesContext";
import "../style/favorites.css";
import ProductCard from "../components/ProductCard";
import { Link } from "react-router-dom";
import { api } from "../api/client";
import { Product } from "../api/types";

const Favorites: React.FC = () => {
  const { favoriteIds, clearFavorites } = useFavorites();
  const [products, setProducts] = useState<Product[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (favoriteIds.size === 0) {
      setProducts([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    Promise.all(
      Array.from(favoriteIds).map(id => api.getProduct(id).catch(() => null))
    ).then(results => {
      setProducts(results.filter(Boolean) as Product[]);
    }).finally(() => setLoading(false));
  }, [favoriteIds]);

  const isEmpty = favoriteIds.size === 0;

  return (
    <div className="favorites-container">
      <h1 className="h-home">Избранные товары</h1>
      {!isEmpty && (
        <div className="favorites-header">
          <span></span>
          <button className="favorites-clear" onClick={clearFavorites}>Очистить</button>
        </div>
      )}

      {loading && <p>Загрузка...</p>}

      <div className="favorites-list">
        {products.map(item => (
          <ProductCard
            key={item.id}
            product={item}
          />
        ))}
      </div>

      {!isEmpty && !loading && (
        <div style={{ width: "90%", maxWidth: 400, margin: "24px auto 0 auto", textAlign: "center" }}>
          <Link to="/catalog" className="continue-shopping">
            Продолжить покупки
          </Link>
        </div>
      )}

      {isEmpty && !loading && (
        <div className="favorites-empty">
          <p className="favorites-empty-text">Тут пока ещё ничего нет...</p>
        </div>
      )}
    </div>
  );
};

export default Favorites;
