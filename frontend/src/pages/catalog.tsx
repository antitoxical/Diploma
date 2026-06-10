import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { api } from "../api/client";
import { Category, PLACEHOLDER_IMAGE, rewriteExternalImageUrl } from "../api/types";
import "../style/catalog.css";

const Catalog: React.FC = () => {
  const [categories, setCategories] = useState<Category[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api.getRootCategories()
      .then(data => setCategories(Array.isArray(data) ? data : []))
      .catch(console.error)
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="catalog-container">
      <h1 className="h-home">Каталог</h1>
      {loading && <p>Загрузка...</p>}
      <div className="catalog-grid">
        {categories.map(cat => (
          <Link to={`/category/${cat.id}`} key={cat.id}>
            <img
              src={cat.image_url ? rewriteExternalImageUrl(cat.image_url) : PLACEHOLDER_IMAGE}
              alt={cat.name}
            />
            <span>{cat.name}</span>
          </Link>
        ))}
      </div>
    </div>
  );
};

export default Catalog;
