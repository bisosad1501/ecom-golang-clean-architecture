-- Remove business constraints

-- Drop constraints
ALTER TABLE carts DROP CONSTRAINT IF EXISTS check_cart_ownership;
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS check_cart_item_quantity;
ALTER TABLE cart_items DROP CONSTRAINT IF EXISTS check_cart_item_price;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_product_price;
ALTER TABLE products DROP CONSTRAINT IF EXISTS check_product_stock;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS check_order_total;
ALTER TABLE orders DROP CONSTRAINT IF EXISTS check_order_subtotal;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS check_order_item_quantity;
ALTER TABLE order_items DROP CONSTRAINT IF EXISTS check_order_item_price;
ALTER TABLE stock_reservations DROP CONSTRAINT IF EXISTS check_reservation_quantity;
ALTER TABLE payments DROP CONSTRAINT IF EXISTS check_payment_amount;
ALTER TABLE reviews DROP CONSTRAINT IF EXISTS check_review_rating;
ALTER TABLE coupons DROP CONSTRAINT IF EXISTS check_coupon_discount_percentage;
ALTER TABLE coupons DROP CONSTRAINT IF EXISTS check_coupon_discount_amount;

-- Drop indexes
DROP INDEX IF EXISTS idx_unique_active_cart_per_user;
DROP INDEX IF EXISTS idx_unique_active_cart_per_session;
DROP INDEX IF EXISTS idx_unique_cart_item_per_product;
DROP INDEX IF EXISTS idx_stock_reservations_cleanup;
DROP INDEX IF EXISTS idx_carts_expiration;
DROP INDEX IF EXISTS idx_orders_payment_timeout;
