-- Add delivery_method_id, payment_method_id, promocode_id, address_snapshot to orders
-- These columns were present in the target schema but missing from the initial DB migration.

ALTER TABLE orders
    ADD COLUMN IF NOT EXISTS delivery_method_id  UUID        REFERENCES delivery_methods (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS payment_method_id   UUID        REFERENCES payment_methods  (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS promocode_id        UUID        REFERENCES promocodes        (id) ON DELETE SET NULL,
    ADD COLUMN IF NOT EXISTS address_snapshot    TEXT;
