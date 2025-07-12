-- Migration script to add new fields for Phase 2 fixes
-- Run this after docker-compose up to add missing fields

-- Add new fields to carts table
ALTER TABLE carts ADD COLUMN IF NOT EXISTS session_id TEXT;
ALTER TABLE carts ADD COLUMN IF NOT EXISTS status TEXT DEFAULT 'active';
ALTER TABLE carts ADD COLUMN IF NOT EXISTS expires_at TIMESTAMP WITH TIME ZONE;
ALTER TABLE carts ADD COLUMN IF NOT EXISTS subtotal NUMERIC DEFAULT 0;
ALTER TABLE carts ADD COLUMN IF NOT EXISTS total NUMERIC DEFAULT 0;
ALTER TABLE carts ADD COLUMN IF NOT EXISTS item_count INTEGER DEFAULT 0;
ALTER TABLE carts ADD COLUMN IF NOT EXISTS currency TEXT DEFAULT 'USD';
ALTER TABLE carts ADD COLUMN IF NOT EXISTS notes TEXT;

-- Fix stock_reservations table for guest cart support
ALTER TABLE stock_reservations ADD COLUMN IF NOT EXISTS session_id TEXT;
ALTER TABLE stock_reservations ALTER COLUMN user_id DROP NOT NULL;

-- Add new fields to payments table
ALTER TABLE payments ADD COLUMN IF NOT EXISTS payment_intent_id TEXT;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS gateway TEXT DEFAULT 'stripe';
ALTER TABLE payments ADD COLUMN IF NOT EXISTS processing_fee NUMERIC DEFAULT 0;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS gateway_fee NUMERIC DEFAULT 0;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS net_amount NUMERIC DEFAULT 0;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS failure_code TEXT;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS refund_reason TEXT;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS metadata TEXT;
ALTER TABLE payments ADD COLUMN IF NOT EXISTS notes TEXT;

-- Add new fields to orders table
ALTER TABLE orders ADD COLUMN IF NOT EXISTS version INTEGER DEFAULT 1;
ALTER TABLE orders ADD COLUMN IF NOT EXISTS last_modified_by UUID;

-- Create new indexes for performance
CREATE INDEX IF NOT EXISTS idx_carts_session_id ON carts(session_id);
CREATE INDEX IF NOT EXISTS idx_carts_status ON carts(status);
CREATE INDEX IF NOT EXISTS idx_carts_expires_at ON carts(expires_at);

CREATE INDEX IF NOT EXISTS idx_payments_payment_intent_id ON payments(payment_intent_id);
CREATE INDEX IF NOT EXISTS idx_payments_gateway ON payments(gateway);

CREATE INDEX IF NOT EXISTS idx_orders_version ON orders(version);
CREATE INDEX IF NOT EXISTS idx_orders_last_modified_by ON orders(last_modified_by);

-- Add missing indexes for stock_reservations
CREATE INDEX IF NOT EXISTS idx_stock_reservations_session_id ON stock_reservations(session_id);

-- Update existing carts to have default values
UPDATE carts SET 
    status = 'active',
    currency = 'USD',
    subtotal = 0,
    total = 0,
    item_count = 0
WHERE status IS NULL;

-- Update existing payments to have default values  
UPDATE payments SET
    gateway = 'stripe',
    processing_fee = 0,
    gateway_fee = 0,
    net_amount = amount
WHERE gateway IS NULL;

-- Update existing orders to have default values
UPDATE orders SET
    version = 1
WHERE version IS NULL;

-- Make user_id nullable for guest carts
ALTER TABLE carts ALTER COLUMN user_id DROP NOT NULL;

-- Add unique constraint for one cart per user (excluding guest carts)
CREATE UNIQUE INDEX IF NOT EXISTS idx_carts_user_id_unique 
ON carts(user_id) WHERE user_id IS NOT NULL;

-- Add unique constraint for one payment per order
CREATE UNIQUE INDEX IF NOT EXISTS idx_payments_order_id_unique 
ON payments(order_id);

COMMIT;
