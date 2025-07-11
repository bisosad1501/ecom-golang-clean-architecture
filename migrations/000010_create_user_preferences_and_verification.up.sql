-- Create user_preferences table
CREATE TABLE IF NOT EXISTS user_preferences (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Display preferences
    theme VARCHAR(20) DEFAULT 'system',
    language VARCHAR(10) DEFAULT 'en',
    currency VARCHAR(10) DEFAULT 'USD',
    timezone VARCHAR(50) DEFAULT 'UTC',
    
    -- Notification preferences
    email_notifications BOOLEAN DEFAULT true,
    sms_notifications BOOLEAN DEFAULT false,
    push_notifications BOOLEAN DEFAULT true,
    marketing_emails BOOLEAN DEFAULT false,
    order_updates BOOLEAN DEFAULT true,
    product_recommendations BOOLEAN DEFAULT false,
    newsletter_subscription BOOLEAN DEFAULT false,
    security_alerts BOOLEAN DEFAULT true,
    
    -- Privacy preferences
    profile_visibility VARCHAR(20) DEFAULT 'private',
    show_online_status BOOLEAN DEFAULT false,
    allow_data_collection BOOLEAN DEFAULT false,
    allow_personalization BOOLEAN DEFAULT true,
    allow_third_party_sharing BOOLEAN DEFAULT false,
    
    -- Shopping preferences
    default_shipping_method VARCHAR(50) DEFAULT '',
    default_payment_method VARCHAR(50) DEFAULT '',
    save_payment_methods BOOLEAN DEFAULT true,
    auto_apply_coupons BOOLEAN DEFAULT true,
    wishlist_visibility VARCHAR(20) DEFAULT 'private',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    UNIQUE(user_id)
);

-- Create user_verifications table
CREATE TABLE IF NOT EXISTS user_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(20) NOT NULL, -- 'email' or 'phone'
    token VARCHAR(255), -- for email verification
    code VARCHAR(10), -- for phone verification
    attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_user_verifications_user_id ON user_verifications(user_id);
CREATE INDEX IF NOT EXISTS idx_user_verifications_token ON user_verifications(token);
CREATE INDEX IF NOT EXISTS idx_user_verifications_code ON user_verifications(code);
CREATE INDEX IF NOT EXISTS idx_user_verifications_type ON user_verifications(type);
CREATE INDEX IF NOT EXISTS idx_user_verifications_expires_at ON user_verifications(expires_at);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_preferences_updated_at 
    BEFORE UPDATE ON user_preferences 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_user_verifications_updated_at
    BEFORE UPDATE ON user_verifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
