import React, { useEffect, useState } from "react";
import { useParams, useSearchParams } from "react-router-dom";
import { api } from "../api/client";
import { Product, Category, Firm, Tag, PaginatedProducts } from "../api/types";
import ProductCard from "../components/ProductCard";
import "../style/catalog.css";
import "../style/home.css";

export const CategoryPage: React.FC = () => {
  const { categoryId } = useParams();
  const [searchParams] = useSearchParams();
  const [products, setProducts] = useState<Product[]>([]);
  const [category, setCategory] = useState<Category | null>(null);
  const [firms, setFirms] = useState<Firm[]>([]);
  const [tags, setTags] = useState<Tag[]>([]);
  const [total, setTotal] = useState(0);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);

  const [selectedFirm, setSelectedFirm] = useState("");
  const [selectedTags, setSelectedTags] = useState<string[]>([]);
  const [sort, setSort] = useState("newest");
  const [inStockOnly, setInStockOnly] = useState(false);
  const [minPrice, setMinPrice] = useState("");
  const [maxPrice, setMaxPrice] = useState("");
  const [minRating, setMinRating] = useState("");
  const [onlyNew, setOnlyNew] = useState(false);
  const [showFilters, setShowFilters] = useState(false);

  const searchQuery = searchParams.get("search") || "";
  const isAll = categoryId === "all";

  useEffect(() => {
    if (!isAll && categoryId) {
      api.getCategory(categoryId).then(setCategory).catch(() => {});
    }
    api.getFirms().then(setFirms).catch(() => {});
    api.getTags().then(setTags).catch(() => {});
  }, [categoryId, isAll]);

  useEffect(() => {
    setLoading(true);
    const params: Record<string, string> = {
      page: String(page),
      limit: "20",
      sort_by: sort,
    };
    if (!isAll && categoryId) {
      params.category_id = categoryId;
      params.include_subcategories = "true";
    }
    if (selectedFirm) params.firm_id = selectedFirm;
    if (selectedTags.length === 1) {
      params.tag_id = selectedTags[0];
    } else if (selectedTags.length > 1) {
      params.tag_ids = selectedTags.join(",");
    }
    if (searchQuery) params.search = searchQuery;
    if (inStockOnly) params.in_stock = "true";
    if (minPrice) params.min_price = String(Number(minPrice) * 100);
    if (maxPrice) params.max_price = String(Number(maxPrice) * 100);
    if (minRating) params.min_rating = minRating;
    if (onlyNew) params.new_days = "14";

    api.searchProducts(params)
      .then((data: PaginatedProducts) => {
        setProducts(data.products || []);
        setTotal(data.total);
        setTotalPages(data.total_pages);
      })
      .catch(console.error)
      .finally(() => setLoading(false));
  }, [categoryId, page, selectedFirm, selectedTags, searchQuery, sort, isAll, inStockOnly, minPrice, maxPrice, minRating, onlyNew]);

  const title = isAll
    ? (searchQuery ? `Поиск: ${searchQuery}` : "Все товары")
    : (category?.name || "Категория");

  const handleTagToggle = (tagId: string) => {
    setSelectedTags(prev =>
      prev.includes(tagId)
        ? prev.filter(id => id !== tagId)
        : [...prev, tagId]
    );
    setPage(1);
  };

  const resetFilters = () => {
    setSelectedFirm("");
    setSelectedTags([]);
    setSort("newest");
    setInStockOnly(false);
    setMinPrice("");
    setMaxPrice("");
    setMinRating("");
    setOnlyNew(false);
    setPage(1);
  };

  return (
    <div className="category-container">
      <h1 className="h-home">{title}</h1>

      <div
        className="category-filters"
        style={{
          display: "flex",
          flexDirection: "column",
          gap: 8,
          width: "92%",
          maxWidth: 500,
          margin: "0 auto 18px",
          boxSizing: "border-box",
        }}
      >
        <div
          className="category-filters-row"
          style={{ display: "flex", gap: 8, width: "100%", minWidth: 0 }}
        >
          <select
            className="sort-select"
            style={{ flex: "1 1 0", minWidth: 0, width: 0 }}
            value={selectedFirm}
            onChange={e => { setSelectedFirm(e.target.value); setPage(1); }}
          >
            <option value="">Все бренды</option>
            {firms.map(f => (
              <option key={f.id} value={f.id}>{f.name}</option>
            ))}
          </select>
          <select
            className="sort-select"
            style={{ flex: "1 1 0", minWidth: 0, width: 0 }}
            value={sort}
            onChange={e => { setSort(e.target.value); setPage(1); }}
          >
            <option value="newest">Новинки</option>
            <option value="price_asc">Цена ↑</option>
            <option value="price_desc">Цена ↓</option>
            <option value="rating">По рейтингу</option>
            <option value="popular">По популярности</option>
          </select>
        </div>
        <button
          type="button"
          className="filter-toggle-btn"
          style={{ width: "100%", maxWidth: "100%", flexShrink: 0 }}
          onClick={() => setShowFilters(p => !p)}
          aria-expanded={showFilters}
        >
          {showFilters ? "Скрыть" : "Фильтры"}
        </button>
      </div>

      {showFilters && (
        <div className="extended-filters">
          {tags.length > 0 && (
            <div className="filter-group">
              <span className="filter-group-label">Тип товара (теги):</span>
              <div className="filter-tags">
                {tags.map(t => (
                  <button
                    key={t.id}
                    className={`filter-tag-btn${selectedTags.includes(t.id) ? " active" : ""}`}
                    onClick={() => handleTagToggle(t.id)}
                    style={selectedTags.includes(t.id) ? { borderColor: t.color || "#a349a4", background: (t.color || "#a349a4") + "22" } : {}}
                  >
                    {t.name}
                  </button>
                ))}
              </div>
            </div>
          )}

          <div className="filter-group">
            <span className="filter-group-label">Цена (Br):</span>
            <div className="filter-price-row">
              <input
                type="number"
                placeholder="От"
                value={minPrice}
                onChange={e => { setMinPrice(e.target.value); setPage(1); }}
                className="filter-price-input"
              />
              <span>—</span>
              <input
                type="number"
                placeholder="До"
                value={maxPrice}
                onChange={e => { setMaxPrice(e.target.value); setPage(1); }}
                className="filter-price-input"
              />
            </div>
          </div>

          <div className="filter-group">
            <span className="filter-group-label">Мин. рейтинг:</span>
            <select
              className="sort-select"
              value={minRating}
              onChange={e => { setMinRating(e.target.value); setPage(1); }}
            >
              <option value="">Любой</option>
              <option value="4">★ 4+</option>
              <option value="3">★ 3+</option>
              <option value="2">★ 2+</option>
            </select>
          </div>

          <div className="filter-group filter-checks">
            <label className="filter-check">
              <input
                type="checkbox"
                checked={inStockOnly}
                onChange={e => { setInStockOnly(e.target.checked); setPage(1); }}
              />
              Только в наличии
            </label>
            <label className="filter-check">
              <input
                type="checkbox"
                checked={onlyNew}
                onChange={e => { setOnlyNew(e.target.checked); setPage(1); }}
              />
              Новинки (14 дней)
            </label>
          </div>

          <button className="filter-reset-btn" onClick={resetFilters}>
            Сбросить фильтры
          </button>
        </div>
      )}

      {loading && <p>Загрузка...</p>}

      <div className="product-list">
        {products.map(item => (
          <ProductCard key={item.id} product={item} />
        ))}
        {!loading && products.length === 0 && <p>Товары не найдены</p>}
      </div>

      {totalPages > 1 && (
        <div className="pagination">
          <button disabled={page <= 1} onClick={() => setPage(p => p - 1)}>←</button>
          <span>{page} / {totalPages}</span>
          <button disabled={page >= totalPages} onClick={() => setPage(p => p + 1)}>→</button>
        </div>
      )}
    </div>
  );
};
