import React, { useEffect, useState, useRef } from "react";
import { useNavigate } from "react-router-dom";
import { api } from "../api/client";
import { Product, getMainImage } from "../api/types";
import ProductCard from "../components/ProductCard";
import "../style/home.css";

const Home: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [search, setSearch] = useState("");
  const [suggestions, setSuggestions] = useState<Product[]>([]);
  const [showSuggestions, setShowSuggestions] = useState(false);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();
  const debounceRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    api.getProducts()
      .then(data => setProducts(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  const handleSearchChange = (val: string) => {
    setSearch(val);
    if (debounceRef.current !== null) clearTimeout(debounceRef.current);

    if (val.trim().length < 2) {
      setSuggestions([]);
      setShowSuggestions(false);
      return;
    }

    debounceRef.current = setTimeout(() => {
      api.searchProducts({ search: val.trim(), limit: "5" })
        .then(data => {
          const items = data?.products || [];
          setSuggestions(items);
          setShowSuggestions(items.length > 0);
        })
        .catch(() => setSuggestions([]));
    }, 300);
  };

  const handleSearch = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter" && search.trim()) {
      setShowSuggestions(false);
      navigate(`/category/all?search=${encodeURIComponent(search.trim())}`);
    }
  };

  return (
    <div className="home-container">
      <h1 className="h-home">Главная</h1>
      <div className="search-wrapper">
        <input
          type="text"
          placeholder="Поиск товаров"
          className="search-input"
          value={search}
          onChange={e => handleSearchChange(e.target.value)}
          onKeyDown={handleSearch}
          onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
          onFocus={() => suggestions.length > 0 && setShowSuggestions(true)}
        />
        {showSuggestions && (
          <div className="search-suggestions">
            {suggestions.map(p => (
              <div
                key={p.id}
                className="search-suggestion-item"
                onMouseDown={() => {
                  setShowSuggestions(false);
                  navigate(`/product/${p.id}`);
                }}
              >
                <img src={getMainImage(p)} alt="" className="suggestion-img" />
                <span>{p.name}</span>
              </div>
            ))}
          </div>
        )}
      </div>
      <div className="product-list">
        <p className="section-title">Выбор покупателей</p>
        {loading && <p>Загрузка...</p>}
        {products.slice(0, 6).map(p => (
          <ProductCard
            key={p.id}
            product={p}
          />
        ))}
      </div>
    </div>
  );
};

export default Home;
