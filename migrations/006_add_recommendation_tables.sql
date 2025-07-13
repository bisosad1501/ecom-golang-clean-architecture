-- Migration: Add recommendation system tables
-- Version: 006
-- Description: Create tables for product recommendations, user interactions, and related features

-- Create user_product_interactions table
CREATE TABLE IF NOT EXISTS user_product_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    session_id VARCHAR(255),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    interaction_type VARCHAR(50) NOT NULL,
    value DECIMAL(10,2) DEFAULT 0,
    metadata TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create product_recommendations table
CREATE TABLE IF NOT EXISTS product_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    recommended_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    recommendation_type VARCHAR(50) NOT NULL,
    score DECIMAL(5,4) DEFAULT 0,
    reason TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, recommended_product_id, recommendation_type)
);

-- Create product_similarities table
CREATE TABLE IF NOT EXISTS product_similarities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    similar_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    similarity_score DECIMAL(5,4) NOT NULL DEFAULT 0,
    similarity_type VARCHAR(50) NOT NULL DEFAULT 'content_based',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, similar_product_id, similarity_type)
);

-- Create frequently_bought_together table
CREATE TABLE IF NOT EXISTS frequently_bought_together (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    bought_with_product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    frequency_score DECIMAL(5,4) NOT NULL DEFAULT 0,
    support_count INTEGER DEFAULT 0,
    confidence DECIMAL(5,4) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, bought_with_product_id)
);

-- Create trending_products table
CREATE TABLE IF NOT EXISTS trending_products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    period VARCHAR(20) NOT NULL DEFAULT 'weekly',
    trending_score DECIMAL(10,4) NOT NULL DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    purchase_count INTEGER DEFAULT 0,
    search_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(product_id, period)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_user_product_interactions_user_id ON user_product_interactions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_product_interactions_product_id ON user_product_interactions(product_id);
CREATE INDEX IF NOT EXISTS idx_user_product_interactions_session_id ON user_product_interactions(session_id);
CREATE INDEX IF NOT EXISTS idx_user_product_interactions_type ON user_product_interactions(interaction_type);
CREATE INDEX IF NOT EXISTS idx_user_product_interactions_created_at ON user_product_interactions(created_at);

CREATE INDEX IF NOT EXISTS idx_product_recommendations_product_id ON product_recommendations(product_id);
CREATE INDEX IF NOT EXISTS idx_product_recommendations_recommended_id ON product_recommendations(recommended_product_id);
CREATE INDEX IF NOT EXISTS idx_product_recommendations_type ON product_recommendations(recommendation_type);
CREATE INDEX IF NOT EXISTS idx_product_recommendations_score ON product_recommendations(score DESC);

CREATE INDEX IF NOT EXISTS idx_product_similarities_product_id ON product_similarities(product_id);
CREATE INDEX IF NOT EXISTS idx_product_similarities_similar_id ON product_similarities(similar_product_id);
CREATE INDEX IF NOT EXISTS idx_product_similarities_score ON product_similarities(similarity_score DESC);

CREATE INDEX IF NOT EXISTS idx_frequently_bought_product_id ON frequently_bought_together(product_id);
CREATE INDEX IF NOT EXISTS idx_frequently_bought_with_id ON frequently_bought_together(bought_with_product_id);
CREATE INDEX IF NOT EXISTS idx_frequently_bought_score ON frequently_bought_together(frequency_score DESC);

CREATE INDEX IF NOT EXISTS idx_trending_products_product_id ON trending_products(product_id);
CREATE INDEX IF NOT EXISTS idx_trending_products_period ON trending_products(period);
CREATE INDEX IF NOT EXISTS idx_trending_products_score ON trending_products(trending_score DESC);
CREATE INDEX IF NOT EXISTS idx_trending_products_updated_at ON trending_products(updated_at);

-- Add some sample data for testing
INSERT INTO trending_products (product_id, period, trending_score, view_count, purchase_count, search_count)
SELECT 
    p.id,
    'weekly',
    RANDOM() * 100,
    FLOOR(RANDOM() * 1000)::INTEGER,
    FLOOR(RANDOM() * 100)::INTEGER,
    FLOOR(RANDOM() * 500)::INTEGER
FROM products p
LIMIT 20
ON CONFLICT (product_id, period) DO NOTHING;

INSERT INTO trending_products (product_id, period, trending_score, view_count, purchase_count, search_count)
SELECT 
    p.id,
    'daily',
    RANDOM() * 100,
    FLOOR(RANDOM() * 500)::INTEGER,
    FLOOR(RANDOM() * 50)::INTEGER,
    FLOOR(RANDOM() * 250)::INTEGER
FROM products p
LIMIT 15
ON CONFLICT (product_id, period) DO NOTHING;

INSERT INTO trending_products (product_id, period, trending_score, view_count, purchase_count, search_count)
SELECT 
    p.id,
    'monthly',
    RANDOM() * 100,
    FLOOR(RANDOM() * 2000)::INTEGER,
    FLOOR(RANDOM() * 200)::INTEGER,
    FLOOR(RANDOM() * 1000)::INTEGER
FROM products p
LIMIT 25
ON CONFLICT (product_id, period) DO NOTHING;
