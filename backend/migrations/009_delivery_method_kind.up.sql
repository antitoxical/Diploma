-- Тип способа доставки задаёт поля на оформлении заказа (курьер / самовывоз / ПВЗ / почта).
ALTER TABLE delivery_methods
    ADD COLUMN IF NOT EXISTS kind VARCHAR(20) NOT NULL DEFAULT 'courier';

ALTER TABLE delivery_methods
    DROP CONSTRAINT IF EXISTS delivery_methods_kind_check;

ALTER TABLE delivery_methods
    ADD CONSTRAINT delivery_methods_kind_check
        CHECK (kind IN ('courier', 'pickup', 'pvz', 'postal'));

UPDATE delivery_methods
SET kind = 'pickup'
WHERE name ILIKE '%самовывоз%'
   OR COALESCE(description, '') ILIKE '%самовывоз%';

UPDATE delivery_methods
SET kind = 'pvz'
WHERE kind = 'courier'
  AND (
    name ILIKE '%сдэк%'
        OR name ILIKE '%cdek%'
        OR COALESCE(description, '') ILIKE '%сдэк%'
        OR COALESCE(description, '') ILIKE '%cdek%'
        OR COALESCE(description, '') ILIKE '%пвз%'
    );

UPDATE delivery_methods
SET kind = 'postal'
WHERE kind = 'courier'
  AND (
    name ILIKE '%почт%'
        OR COALESCE(description, '') ILIKE '%почт%'
    );
