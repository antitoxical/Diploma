import React, { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { api } from "../api/client";
import { Review } from "../api/types";
import { useTelegram } from "../context/TelegramContext";
import { useToast } from "../context/ToastContext";
import "../style/reviews.css";

const Reviews: React.FC = () => {
  const { productId } = useParams();
  const navigate = useNavigate();
  const { user: tgUser } = useTelegram();
  const { showToast } = useToast();
  const [reviews, setReviews] = useState<Review[]>([]);
  const [loading, setLoading] = useState(true);
  const [newRating, setNewRating] = useState(5);
  const [newComment, setNewComment] = useState("");
  const [submitting, setSubmitting] = useState(false);

  const loadReviews = () => {
    if (!productId) return;
    setLoading(true);
    api.getProductReviews(productId)
      .then(data => setReviews(Array.isArray(data) ? data : []))
      .catch(() => setReviews([]))
      .finally(() => setLoading(false));
  };

  useEffect(() => { loadReviews(); }, [productId]);

  const handleSubmit = async () => {
    if (!tgUser || !productId) {
      showToast("Откройте через Telegram для отправки отзыва", "error");
      return;
    }
    setSubmitting(true);
    try {
      await api.upsertReview({
        user_id: tgUser.id,
        product_id: productId,
        rating: newRating,
        comment: newComment || undefined,
      });
      showToast("Отзыв отправлен!");
      setNewComment("");
      setNewRating(5);
      loadReviews();
    } catch {
      showToast("Ошибка отправки", "error");
    }
    setSubmitting(false);
  };

  return (
    <div className="reviews-container">
      <div className="reviews-header">
        <button className="reviews-back" onClick={() => navigate(-1)}>&lt;</button>
        <h2 className="reviews-title">Отзывы</h2>
      </div>

      {loading && <p className="reviews-loading">Загрузка...</p>}

      <div className="reviews-list">
        {!loading && reviews.length === 0 && (
          <p className="reviews-empty">Отзывов пока нет. Будьте первым!</p>
        )}
        {reviews.map((review) => (
          <div className="review-card" key={review.id}>
            <div className="review-header">
              <span className="review-user">Пользователь</span>
              <span className="review-stars">
                {Array.from({ length: 5 }).map((_, idx) =>
                  idx < review.rating ? (
                    <span key={idx} className="star filled">★</span>
                  ) : (
                    <span key={idx} className="star">☆</span>
                  )
                )}
              </span>
            </div>
            {review.comment && <div className="review-text">{review.comment}</div>}
          </div>
        ))}
      </div>

      <div className="review-form">
        <h3 className="review-form-title">Оставить отзыв</h3>
        <div className="review-form-rating">
          {Array.from({ length: 5 }).map((_, idx) => (
            <button
              key={idx}
              className={`review-star-btn ${idx < newRating ? "active" : ""}`}
              onClick={() => setNewRating(idx + 1)}
            >
              {idx < newRating ? "★" : "☆"}
            </button>
          ))}
        </div>
        <textarea
          className="review-form-comment"
          placeholder="Ваш комментарий (необязательно)"
          value={newComment}
          onChange={e => setNewComment(e.target.value)}
          rows={3}
        />
        <button
          className="review-form-submit"
          disabled={submitting}
          onClick={handleSubmit}
        >
          {submitting ? "Отправка..." : "Отправить отзыв"}
        </button>
      </div>
    </div>
  );
};

export default Reviews;
