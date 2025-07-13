-- Create product_recommendations table
CREATE TABLE IF NOT EXISTS product_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    recommended_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    type VARCHAR(50) NOT NULL,
    score DECIMAL(5,4) DEFAULT 0,
    reason TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for product_recommendations
CREATE INDEX idx_product_recommendations_product_id ON product_recommendations(product_id);
CREATE INDEX idx_product_recommendations_recommended_id ON product_recommendations(recommended_id);
CREATE INDEX idx_product_recommendations_type ON product_recommendations(type);
CREATE INDEX idx_product_recommendations_score ON product_recommendations(score DESC);
CREATE UNIQUE INDEX idx_product_recommendations_unique ON product_recommendations(product_id, recommended_id, type);

-- Create user_product_interactions table
CREATE TABLE IF NOT EXISTS user_product_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(128),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    interaction_type VARCHAR(50) NOT NULL,
    value DECIMAL(10,2) DEFAULT 1,
    metadata TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for user_product_interactions
CREATE INDEX idx_user_product_interactions_user_id ON user_product_interactions(user_id);
CREATE INDEX idx_user_product_interactions_session_id ON user_product_interactions(session_id);
CREATE INDEX idx_user_product_interactions_product_id ON user_product_interactions(product_id);
CREATE INDEX idx_user_product_interactions_type ON user_product_interactions(interaction_type);
CREATE INDEX idx_user_product_interactions_created_at ON user_product_interactions(created_at);

-- Create product_similarities table
CREATE TABLE IF NOT EXISTS product_similarities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    similar_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    similarity_score DECIMAL(5,4) NOT NULL,
    algorithm VARCHAR(50) NOT NULL,
    features TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for product_similarities
CREATE INDEX idx_product_similarities_product_id ON product_similarities(product_id);
CREATE INDEX idx_product_similarities_similar_id ON product_similarities(similar_id);
CREATE INDEX idx_product_similarities_score ON product_similarities(similarity_score DESC);
CREATE UNIQUE INDEX idx_product_similarities_unique ON product_similarities(product_id, similar_id, algorithm);

-- Create frequently_bought_together table
CREATE TABLE IF NOT EXISTS frequently_bought_together (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    with_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    frequency INTEGER DEFAULT 1,
    confidence DECIMAL(5,4) DEFAULT 0,
    support DECIMAL(5,4) DEFAULT 0,
    lift DECIMAL(8,4) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for frequently_bought_together
CREATE INDEX idx_frequently_bought_together_product_id ON frequently_bought_together(product_id);
CREATE INDEX idx_frequently_bought_together_with_id ON frequently_bought_together(with_id);
CREATE INDEX idx_frequently_bought_together_confidence ON frequently_bought_together(confidence DESC);
CREATE INDEX idx_frequently_bought_together_frequency ON frequently_bought_together(frequency DESC);
CREATE UNIQUE INDEX idx_frequently_bought_together_unique ON frequently_bought_together(product_id, with_id);

-- Create trending_products table
CREATE TABLE IF NOT EXISTS trending_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    trend_score DECIMAL(10,2) NOT NULL,
    view_count INTEGER DEFAULT 0,
    sales_count INTEGER DEFAULT 0,
    search_count INTEGER DEFAULT 0,
    period VARCHAR(20) NOT NULL,
    date DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for trending_products
CREATE INDEX idx_trending_products_product_id ON trending_products(product_id);
CREATE INDEX idx_trending_products_trend_score ON trending_products(trend_score DESC);
CREATE INDEX idx_trending_products_period ON trending_products(period);
CREATE INDEX idx_trending_products_date ON trending_products(date);
CREATE UNIQUE INDEX idx_trending_products_unique ON trending_products(product_id, period, date);

-- Add constraints to ensure data integrity
ALTER TABLE user_product_interactions 
ADD CONSTRAINT chk_user_or_session CHECK (
    (user_id IS NOT NULL AND session_id IS NULL) OR 
    (user_id IS NULL AND session_id IS NOT NULL)
);

-- Add check constraints for valid interaction types
ALTER TABLE user_product_interactions 
ADD CONSTRAINT chk_interaction_type CHECK (
    interaction_type IN ('view', 'add_to_cart', 'remove_from_cart', 'purchase', 'wishlist', 'review', 'share', 'compare', 'search', 'click')
);

-- Add check constraints for valid recommendation types
ALTER TABLE product_recommendations 
ADD CONSTRAINT chk_recommendation_type CHECK (
    type IN ('related', 'similar', 'frequently_bought', 'trending', 'personalized', 'cross_sell', 'up_sell', 'recently_viewed', 'based_on_category', 'based_on_brand')
);

-- Add check constraints for valid periods
ALTER TABLE trending_products 
ADD CONSTRAINT chk_period CHECK (
    period IN ('daily', 'weekly', 'monthly')
);

-- Ensure products cannot be recommended to themselves
ALTER TABLE product_recommendations 
ADD CONSTRAINT chk_no_self_recommendation CHECK (product_id != recommended_id);

ALTER TABLE product_similarities 
ADD CONSTRAINT chk_no_self_similarity CHECK (product_id != similar_id);

ALTER TABLE frequently_bought_together 
ADD CONSTRAINT chk_no_self_bought_together CHECK (product_id != with_id);

-- Add triggers to update updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_product_recommendations_updated_at 
    BEFORE UPDATE ON product_recommendations 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_product_similarities_updated_at 
    BEFORE UPDATE ON product_similarities 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_frequently_bought_together_updated_at 
    BEFORE UPDATE ON frequently_bought_together 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
