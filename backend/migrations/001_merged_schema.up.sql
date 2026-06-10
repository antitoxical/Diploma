-- Полная схема TelegramShop (26 таблиц + 1 VIEW)
-- PostgreSQL 12+

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ==========================================
-- БЛОК 1: ПОЛЬЗОВАТЕЛИ И БЕЗОПАСНОСТЬ
-- ==========================================

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    telegram_id BIGINT NOT NULL UNIQUE,
    username VARCHAR(255),
    first_name VARCHAR(255),
    last_name VARCHAR(255),
    phone VARCHAR(64),
    delivery_address TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE user_addresses (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    label VARCHAR(255),
    full_address TEXT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE admin_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID NOT NULL REFERENCES admins (id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    ip_address VARCHAR(45),
    user_agent TEXT,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    admin_id UUID REFERENCES admins (id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    table_name VARCHAR(100),
    record_id UUID,
    details JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==========================================
-- БЛОК 2: КАТАЛОГ
-- ==========================================

CREATE TABLE firms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    parent_id UUID REFERENCES categories (id) ON DELETE SET NULL,
    image_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_categories_parent_id ON categories (parent_id);

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    category_id UUID NOT NULL REFERENCES categories (id) ON DELETE RESTRICT,
    firm_id UUID REFERENCES firms (id) ON DELETE SET NULL,
    name VARCHAR(500) NOT NULL,
    description TEXT,
    price_kopeks BIGINT NOT NULL CHECK (price_kopeks >= 0),
    stock_quantity INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_main BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0
);

CREATE TABLE product_price_tiers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    min_quantity INT NOT NULL DEFAULT 1 CHECK (min_quantity >= 1),
    price_kopeks BIGINT NOT NULL CHECK (price_kopeks >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (product_id, min_quantity)
);

CREATE TABLE tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    color VARCHAR(7),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE product_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES tags (id) ON DELETE CASCADE,
    UNIQUE (product_id, tag_id)
);

-- ==========================================
-- БЛОК 3: КОРЗИНА И ПРОЦЕСС ПОКУПКИ
-- ==========================================

CREATE TABLE basket_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    quantity INT NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, product_id)
);

CREATE TABLE delivery_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    price_kopeks BIGINT NOT NULL DEFAULT 0 CHECK (price_kopeks >= 0),
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE payment_methods (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE promocodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) NOT NULL UNIQUE,
    discount_percent INT NOT NULL CHECK (discount_percent > 0 AND discount_percent <= 100),
    max_uses INT,
    current_uses INT NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==========================================
-- БЛОК 4: ЗАКАЗЫ И ТРАНЗАКЦИИ
-- ==========================================

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    delivery_method_id UUID REFERENCES delivery_methods (id) ON DELETE SET NULL,
    payment_method_id UUID REFERENCES payment_methods (id) ON DELETE SET NULL,
    promocode_id UUID REFERENCES promocodes (id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL,
    total_amount_kopeks BIGINT NOT NULL CHECK (total_amount_kopeks >= 0),
    address_snapshot TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT orders_status_check CHECK (
        status IN (
            'created',
            'paid',
            'delivering',
            'completed',
            'cancelled'
        )
    )
);

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE RESTRICT,
    quantity INT NOT NULL CHECK (quantity > 0),
    price_at_purchase_kopeks BIGINT NOT NULL CHECK (price_at_purchase_kopeks >= 0)
);

CREATE TABLE order_status_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    comment TEXT,
    changed_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE order_promocodes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    promocode_id UUID NOT NULL REFERENCES promocodes (id) ON DELETE RESTRICT,
    discount_percent INT NOT NULL,
    discount_amount_kopeks BIGINT NOT NULL DEFAULT 0,
    UNIQUE (order_id)
);

-- ==========================================
-- БЛОК 5: ВЗАИМОДЕЙСТВИЕ
-- ==========================================

CREATE TABLE favorites (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, product_id)
);

CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products (id) ON DELETE CASCADE,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    is_moderated BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (user_id, product_id)
);

CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    title VARCHAR(500),
    message TEXT,
    is_read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE search_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    query VARCHAR(500) NOT NULL,
    results_count INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- ==========================================
-- ИНДЕКСЫ
-- ==========================================

CREATE INDEX idx_user_addresses_user_id ON user_addresses (user_id);
CREATE INDEX idx_admin_sessions_admin_id ON admin_sessions (admin_id);
CREATE INDEX idx_admin_sessions_expires ON admin_sessions (expires_at);
CREATE INDEX idx_audit_logs_admin_id ON audit_logs (admin_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs (created_at);

CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_firm_id ON products (firm_id);
CREATE INDEX idx_products_active ON products (is_active) WHERE is_active = true;
CREATE INDEX idx_product_images_product_id ON product_images (product_id);
CREATE INDEX idx_product_tags_product_id ON product_tags (product_id);
CREATE INDEX idx_product_tags_tag_id ON product_tags (tag_id);
CREATE INDEX idx_product_price_tiers_product_id ON product_price_tiers (product_id);

CREATE INDEX idx_basket_items_user_id ON basket_items (user_id);

CREATE INDEX idx_orders_user_id ON orders (user_id);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_order_items_order_id ON order_items (order_id);
CREATE INDEX idx_order_status_history_order_id ON order_status_history (order_id);
CREATE INDEX idx_order_promocodes_order_id ON order_promocodes (order_id);

CREATE INDEX idx_favorites_user_id ON favorites (user_id);
CREATE INDEX idx_reviews_product_id ON reviews (product_id);
CREATE INDEX idx_notifications_user_id ON notifications (user_id);
CREATE INDEX idx_notifications_unread ON notifications (user_id, is_read)
    WHERE is_read = false;
CREATE INDEX idx_search_history_user_id ON search_history (user_id);
CREATE INDEX idx_search_history_created_at ON search_history (created_at);

-- ==========================================
-- ПРЕДСТАВЛЕНИЯ
-- ==========================================

CREATE OR REPLACE VIEW product_rating_stats AS
SELECT
    product_id,
    ROUND(AVG(rating)::numeric, 2) AS avg_rating,
    COUNT(*)::int AS review_count
FROM reviews
GROUP BY product_id;
