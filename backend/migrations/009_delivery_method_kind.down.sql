ALTER TABLE delivery_methods DROP CONSTRAINT IF EXISTS delivery_methods_kind_check;
ALTER TABLE delivery_methods DROP COLUMN IF EXISTS kind;
