import React, { useEffect, useState } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { api } from "../api/client";
import { Product, ProductRatingStat, formatPrice, getMainImage, getAllImages } from "../api/types";
import { useCart } from "../context/CartContext";
import { useFavorites } from "../context/FavoritesContext";
import { useToast } from "../context/ToastContext";
import "../style/product.css";

const ProductPage: React.FC = () => {
  const { productId } = useParams();
  const navigate = useNavigate();
  const { addToCart } = useCart();
  const { isFavorite, toggleFavorite } = useFavorites();
  const { showToast } = useToast();

  const [product, setProduct] = useState<Product | null>(null);
  const [rating, setRating] = useState<ProductRatingStat | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentImg, setCurrentImg] = useState(0);

  useEffect(() => {
    if (!productId) return;
    setLoading(true);
    Promise.all([
      api.getProduct(productId),
      api.getProductRating(productId).catch(() => null),
    ])
      .then(([p, r]) => { setProduct(p); setRating(r); })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [productId]);

  if (loading) return <div className="product-container"><p>Загрузка...</p></div>;
  if (!product) return <div className="product-container"><p>Товар не найден</p></div>;

  const images = getAllImages(product);
  const fav = isFavorite(product.id);
  const inStock = product.stock_quantity > 0;
  const avgRating = rating?.avg_rating ?? 0;
  const reviewCount = rating?.review_count ?? 0;

  return (
    <div className="product-container">
      <div className="product-header">
        <button className="back-btn" onClick={() => navigate(-1)}>
          <span className="back-arrow">&lt;</span>
        </button>
        <span className="product-title-top">Магазин</span>
        <button className="close-btn" onClick={() => navigate("/")}>×</button>
      </div>

      <div className="product-image-slider-block">
        <img
          src={images[currentImg]}
          alt={product.name}
          className="product-image-large"
        />
        {images.length > 1 && (
          <div className="product-image-dots">
            {images.map((_, idx) => (
              <span
                key={idx}
                className={`dot${idx === currentImg ? " active" : ""}`}
                onClick={() => setCurrentImg(idx)}
              />
            ))}
          </div>
        )}
      </div>

      <div className="product-title">{product.name}</div>

      {product.tags && product.tags.length > 0 && (
        <div className="product-tags">
          {product.tags.map(tag => (
            <span key={tag.id} className="product-tag" style={{ borderColor: tag.color || "#a349a4" }}>
              {tag.name}
            </span>
          ))}
        </div>
      )}

      <div className="product-rating-row">
        <span className="star">★</span>
        <span className="product-rating-value">{avgRating.toFixed(1)}</span>
        <Link to={`/product/${product.id}/reviews`} className="product-reviews-link">
          {reviewCount > 0 ? `${reviewCount} отзывов` : "Отзывы"}
        </Link>
        <button
          className={`product-fav-btn${fav ? " active" : ""}`}
          onClick={() => {
            toggleFavorite(product.id);
            showToast(fav ? "Удалено из избранного" : "Добавлено в избранное", fav ? undefined : "info");
          }}
          style={{ marginLeft: "12px" }}
        >
          {fav ? "♥" : "♡"}
        </button>
      </div>

      <div className="product-description">
        {product.description || "Описание товара..."}
      </div>

      <div className="product-buy-block">
        <div className="product-price-block">
          <span className="product-price-label">Цена:</span>
          <span className="product-price-value">{formatPrice(product.price_kopeks)}</span>
        </div>
        <button
          className="product-add-to-cart"
          disabled={!inStock}
          onClick={() => {
            const added = addToCart({
              id: product.id,
              title: product.name,
              price: product.price_kopeks,
              image: getMainImage(product),
              stockMax: product.stock_quantity,
            });
            if (added) showToast("Добавлено в корзину");
            else showToast(`Доступно не более ${product.stock_quantity} шт.`, "error");
          }}
        >
          Добавить в корзину
        </button>
      </div>

      <div className="product-stock-row">
        {inStock ? (
          <>
            <span className="stock-icon in-stock">✔</span>
            <span className="stock-text in-stock">В наличии ({product.stock_quantity} шт.)</span>
          </>
        ) : (
          <>
            <span className="stock-icon out-of-stock">✖</span>
            <span className="stock-text out-of-stock">Нет в наличии</span>
          </>
        )}
      </div>
    </div>
  );
};

export default ProductPage;
