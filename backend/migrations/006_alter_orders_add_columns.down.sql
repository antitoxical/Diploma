ALTER TABLE orders
    DROP COLUMN IF EXISTS delivery_method_id,
    DROP COLUMN IF EXISTS payment_method_id,
    DROP COLUMN IF EXISTS promocode_id,
    DROP COLUMN IF EXISTS address_snapshot;
