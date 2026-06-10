import React, { useEffect, useState } from "react";
import { useCart } from "../context/CartContext";
import { useFavorites } from "../context/FavoritesContext";
import { useToast } from "../context/ToastContext";
import { Product, formatPrice, getMainImage, PLACEHOLDER_IMAGE } from "../api/types";
import "../style/home.css";
import "../style/productCard.css";
import { useNavigate } from "react-router-dom";

interface ProductCardProps {
  product: Product;
}

const ProductCard: React.FC<ProductCardProps> = ({ product }) => {
  const { addToCart } = useCart();
  const { isFavorite, toggleFavorite } = useFavorites();
  const { showToast } = useToast();
  const navigate = useNavigate();

  const image = getMainImage(product);
  const [imgSrc, setImgSrc] = useState(image);
  const inStock = product.stock_quantity > 0;
  const fav = isFavorite(product.id);

  useEffect(() => {
    setImgSrc(image);
  }, [image]);

  return (
    <div className="product-card" onClick={() => navigate(`/product/${product.id}`)} style={{ cursor: "pointer" }}>
      <div className="product-image-wrap">
        <img
          src={imgSrc}
          alt={product.name}
          className="product-image"
          onError={() => setImgSrc(PLACEHOLDER_IMAGE)}
        />
      </div>
      <div className="product-info">
        <div className="product-title">{product.name}</div>
        <div className="product-meta">
          <span className={inStock ? "in-stock" : "out-of-stock"}>
            {inStock ? "Есть в наличии" : "Нет в наличии"}
          </span>
        </div>
        <div className="product-price">Цена: {formatPrice(product.price_kopeks)}</div>
        <div className="product-actions">
          <button
            className={`add-to-fav${fav ? " active" : ""}`}
            title={fav ? "Убрать из избранного" : "В избранное"}
            onClick={(e) => {
              e.stopPropagation();
              toggleFavorite(product.id);
              showToast(fav ? "Удалено из избранного" : "Добавлено в избранное", fav ? undefined : "info");
            }}
          >
            {fav ? "♥" : "♡"}
          </button>
          <button
            className="add-to-cart"
            disabled={!inStock}
            onClick={(e) => {
              e.stopPropagation();
              const added = addToCart({
                id: product.id,
                title: product.name,
                price: product.price_kopeks,
                image,
                stockMax: product.stock_quantity,
              });
              if (added) showToast("Добавлено в корзину");
              else showToast(`Доступно не более ${product.stock_quantity} шт.`, "error");
            }}
          >
            В корзину
          </button>
        </div>
      </div>
    </div>
  );
};

export default ProductCard;
