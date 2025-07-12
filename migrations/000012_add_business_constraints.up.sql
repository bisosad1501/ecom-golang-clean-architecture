-- Add business constraints to ensure data integrity

-- 1. Ensure cart has either user_id or session_id, but not both
ALTER TABLE carts ADD CONSTRAINT check_cart_ownership 
CHECK (
    (user_id IS NOT NULL AND session_id IS NULL) OR 
    (user_id IS NULL AND session_id IS NOT NULL)
);

-- 2. Ensure cart item quantity is positive
ALTER TABLE cart_items ADD CONSTRAINT check_cart_item_quantity 
CHECK (quantity > 0);

-- 3. Ensure cart item price is non-negative
ALTER TABLE cart_items ADD CONSTRAINT check_cart_item_price 
CHECK (price >= 0);

-- 4. Ensure product price is positive
ALTER TABLE products ADD CONSTRAINT check_product_price 
CHECK (price > 0);

-- 5. Ensure product stock is non-negative
ALTER TABLE products ADD CONSTRAINT check_product_stock 
CHECK (stock >= 0);

-- 6. Ensure order total is non-negative
ALTER TABLE orders ADD CONSTRAINT check_order_total 
CHECK (total >= 0);

-- 7. Ensure order subtotal is non-negative
ALTER TABLE orders ADD CONSTRAINT check_order_subtotal 
CHECK (subtotal >= 0);

-- 8. Ensure order item quantity is positive
ALTER TABLE order_items ADD CONSTRAINT check_order_item_quantity 
CHECK (quantity > 0);

-- 9. Ensure order item price is non-negative
ALTER TABLE order_items ADD CONSTRAINT check_order_item_price 
CHECK (price >= 0);

-- 10. Ensure stock reservation quantity is positive
ALTER TABLE stock_reservations ADD CONSTRAINT check_reservation_quantity 
CHECK (quantity > 0);

-- 11. Ensure payment amount is positive
ALTER TABLE payments ADD CONSTRAINT check_payment_amount 
CHECK (amount > 0);

-- 12. Add unique constraint for cart per user (only one active cart per user)
CREATE UNIQUE INDEX idx_unique_active_cart_per_user 
ON carts (user_id) 
WHERE status = 'active' AND user_id IS NOT NULL;

-- 13. Add unique constraint for cart per session (only one active cart per session)
CREATE UNIQUE INDEX idx_unique_active_cart_per_session 
ON carts (session_id) 
WHERE status = 'active' AND session_id IS NOT NULL;

-- 14. Add unique constraint for cart item per cart (no duplicate products in same cart)
CREATE UNIQUE INDEX idx_unique_cart_item_per_product 
ON cart_items (cart_id, product_id);

-- 15. Add index for stock reservations cleanup
CREATE INDEX idx_stock_reservations_cleanup 
ON stock_reservations (status, expires_at) 
WHERE status = 'active';

-- 16. Add index for cart expiration cleanup
CREATE INDEX idx_carts_expiration 
ON carts (status, expires_at) 
WHERE expires_at IS NOT NULL;

-- 17. Add index for order payment timeout cleanup
CREATE INDEX idx_orders_payment_timeout 
ON orders (payment_status, payment_timeout) 
WHERE payment_timeout IS NOT NULL;

-- 18. Ensure review rating is between 1 and 5
ALTER TABLE reviews ADD CONSTRAINT check_review_rating 
CHECK (rating >= 1 AND rating <= 5);

-- 19. Ensure coupon discount is between 0 and 100 for percentage discounts
ALTER TABLE coupons ADD CONSTRAINT check_coupon_discount_percentage 
CHECK (
    discount_type != 'percentage' OR 
    (discount_value >= 0 AND discount_value <= 100)
);

-- 20. Ensure coupon discount is positive for fixed amount discounts
ALTER TABLE coupons ADD CONSTRAINT check_coupon_discount_amount 
CHECK (
    discount_type != 'fixed_amount' OR 
    discount_value > 0
);
