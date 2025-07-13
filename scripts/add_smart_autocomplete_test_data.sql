-- Smart Autocomplete Test Data Script
-- This script adds comprehensive test data for smart autocomplete functionality

-- First, ensure we have the required extensions
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Clear existing autocomplete data
DELETE FROM autocomplete_entries;
DELETE FROM autocomplete_analytics;

-- Insert comprehensive autocomplete test data
INSERT INTO autocomplete_entries (
    type, value, display_text, priority, search_count, click_count,
    is_trending, is_personalized, synonyms, tags, score, language, metadata
) VALUES
-- Product suggestions
('product', 'iPhone 15', 'iPhone 15 - Latest Apple Smartphone', 95, 150, 45, true, false,
 ARRAY['phone', 'smartphone', 'apple phone', 'iphone'],
 ARRAY['electronics', 'apple', 'trending', 'popular'],
 92.5, 'en', '{"price": 999, "image": "/images/iphone15.jpg", "category": "Electronics"}'),

('product', 'Samsung Galaxy S24', 'Samsung Galaxy S24 - Android Flagship', 90, 120, 35, true, false,
 ARRAY['galaxy', 'samsung phone', 'android phone', 's24'],
 ARRAY['electronics', 'samsung', 'trending', 'android'],
 88.0, 'en', '{"price": 899, "image": "/images/galaxy-s24.jpg", "category": "Electronics"}'),

('product', 'MacBook Pro', 'MacBook Pro - Professional Laptop', 85, 80, 25, false, false,
 ARRAY['laptop', 'macbook', 'apple laptop', 'pro laptop'],
 ARRAY['computers', 'apple', 'professional'],
 82.0, 'en', '{"price": 1999, "image": "/images/macbook-pro.jpg", "category": "Computers"}'),

('product', 'AirPods Pro', 'AirPods Pro - Wireless Earbuds', 80, 95, 30, true, false,
 ARRAY['earbuds', 'headphones', 'wireless earbuds', 'airpods'],
 ARRAY['audio', 'apple', 'trending'],
 85.5, 'en', '{"price": 249, "image": "/images/airpods-pro.jpg", "category": "Audio"}'),

('product', 'Nike Air Max', 'Nike Air Max - Running Shoes', 75, 60, 20, false, false,
 ARRAY['shoes', 'sneakers', 'running shoes', 'nike shoes'],
 ARRAY['footwear', 'nike', 'sports'],
 78.0, 'en', '{"price": 120, "image": "/images/nike-air-max.jpg", "category": "Footwear"}'),

-- Category suggestions
('category', 'Electronics', 'Electronics & Technology', 90, 200, 60, false, false,
 ARRAY['tech', 'gadgets', 'devices', 'electronic'],
 ARRAY['category', 'popular'],
 89.0, 'en', '{"product_count": 1250, "subcategories": ["Smartphones", "Laptops", "Audio"]}'),

('category', 'Clothing', 'Clothing & Fashion', 85, 180, 55, false, false,
 ARRAY['clothes', 'fashion', 'apparel', 'wear'],
 ARRAY['category', 'popular'],
 86.5, 'en', '{"product_count": 2100, "subcategories": ["Men", "Women", "Kids"]}'),

('category', 'Home & Garden', 'Home & Garden Essentials', 80, 140, 40, false, false,
 ARRAY['home', 'garden', 'house', 'furniture'],
 ARRAY['category', 'lifestyle'],
 81.0, 'en', '{"product_count": 890, "subcategories": ["Furniture", "Decor", "Garden"]}'),

('category', 'Sports & Outdoors', 'Sports & Outdoor Equipment', 75, 110, 35, false, false,
 ARRAY['sports', 'outdoor', 'fitness', 'exercise'],
 ARRAY['category', 'active'],
 77.5, 'en', '{"product_count": 650, "subcategories": ["Fitness", "Outdoor", "Team Sports"]}'),

-- Brand suggestions
('brand', 'Apple', 'Apple Inc. - Premium Technology', 95, 300, 90, false, false,
 ARRAY['apple inc', 'iphone maker', 'mac', 'ios'],
 ARRAY['brand', 'premium', 'technology'],
 94.0, 'en', '{"product_count": 45, "rating": 4.8, "founded": 1976}'),

('brand', 'Samsung', 'Samsung - Innovation & Quality', 90, 250, 75, false, false,
 ARRAY['samsung electronics', 'galaxy maker', 'korean brand'],
 ARRAY['brand', 'electronics', 'innovation'],
 89.5, 'en', '{"product_count": 38, "rating": 4.6, "founded": 1938}'),

('brand', 'Nike', 'Nike - Just Do It', 85, 200, 65, false, false,
 ARRAY['nike inc', 'swoosh', 'athletic wear'],
 ARRAY['brand', 'sports', 'athletic'],
 87.0, 'en', '{"product_count": 120, "rating": 4.7, "founded": 1964}'),

('brand', 'Sony', 'Sony - Entertainment & Electronics', 80, 150, 45, false, false,
 ARRAY['sony corporation', 'playstation maker', 'electronics'],
 ARRAY['brand', 'entertainment', 'electronics'],
 82.5, 'en', '{"product_count": 65, "rating": 4.5, "founded": 1946}'),

-- Query suggestions (popular searches)
('query', 'best smartphone 2024', 'Best Smartphone 2024 Reviews', 85, 180, 50, true, false,
 ARRAY['top phones', 'smartphone reviews', 'best phones'],
 ARRAY['query', 'trending', 'reviews'],
 86.0, 'en', '{"result_count": 245, "avg_price": 750}'),

('query', 'wireless headphones', 'Wireless Headphones & Earbuds', 80, 160, 45, true, false,
 ARRAY['bluetooth headphones', 'earbuds', 'wireless audio'],
 ARRAY['query', 'trending', 'audio'],
 83.5, 'en', '{"result_count": 189, "avg_price": 150}'),

('query', 'gaming laptop', 'Gaming Laptops for Gamers', 75, 120, 35, false, false,
 ARRAY['gaming computer', 'game laptop', 'gaming pc'],
 ARRAY['query', 'gaming', 'computers'],
 78.0, 'en', '{"result_count": 67, "avg_price": 1200}'),

('query', 'running shoes', 'Running Shoes & Athletic Footwear', 70, 100, 30, false, false,
 ARRAY['athletic shoes', 'sport shoes', 'jogging shoes'],
 ARRAY['query', 'sports', 'footwear'],
 75.5, 'en', '{"result_count": 156, "avg_price": 95}'),

-- Trending suggestions
('product', 'Steam Deck', 'Steam Deck - Portable Gaming', 70, 45, 15, true, false,
 ARRAY['handheld gaming', 'portable console', 'steam portable'],
 ARRAY['gaming', 'trending', 'portable'],
 72.0, 'en', '{"price": 399, "image": "/images/steam-deck.jpg", "category": "Gaming"}'),

('product', 'Tesla Model 3', 'Tesla Model 3 - Electric Vehicle', 65, 35, 12, true, false,
 ARRAY['electric car', 'tesla car', 'ev'],
 ARRAY['automotive', 'trending', 'electric'],
 68.0, 'en', '{"price": 35000, "image": "/images/tesla-model3.jpg", "category": "Automotive"}'),

-- Personalized suggestions (these would normally be user-specific)
('query', 'apple accessories', 'Apple Accessories & Cases', 60, 25, 8, false, true,
 ARRAY['iphone case', 'apple charger', 'mac accessories'],
 ARRAY['query', 'personalized', 'accessories'],
 65.0, 'en', '{"result_count": 89, "user_interest": "apple_products"}'),

('product', 'iPad Pro', 'iPad Pro - Professional Tablet', 75, 55, 18, false, true,
 ARRAY['tablet', 'ipad', 'apple tablet'],
 ARRAY['tablets', 'apple', 'personalized'],
 76.5, 'en', '{"price": 799, "image": "/images/ipad-pro.jpg", "category": "Tablets"}');

-- Insert some synonym data for better fuzzy matching
INSERT INTO search_synonyms (term, synonyms, is_active) VALUES
('phone', ARRAY['smartphone', 'mobile', 'cell phone', 'cellular'], true),
('laptop', ARRAY['notebook', 'computer', 'pc', 'portable computer'], true),
('headphones', ARRAY['earphones', 'earbuds', 'audio', 'headset'], true),
('shoes', ARRAY['footwear', 'sneakers', 'boots', 'sandals'], true),
('tv', ARRAY['television', 'smart tv', 'display', 'monitor'], true),
('watch', ARRAY['smartwatch', 'timepiece', 'wearable', 'fitness tracker'], true);

-- Insert some analytics data to simulate user interactions
INSERT INTO autocomplete_analytics (entry_id, interaction_type, query, position, created_at)
SELECT
    ae.id,
    CASE
        WHEN random() < 0.7 THEN 'impression'
        ELSE 'click'
    END,
    'test query',
    floor(random() * 10 + 1)::integer,
    NOW() - (random() * interval '30 days')
FROM autocomplete_entries ae
WHERE random() < 0.3; -- Only for 30% of entries

-- Update trending status based on recent activity
UPDATE autocomplete_entries
SET is_trending = true
WHERE search_count > 100 OR (created_at >= NOW() - INTERVAL '7 days' AND search_count > 50);

-- Calculate initial scores for all entries
UPDATE autocomplete_entries
SET score = (
    (search_count * 0.4) +
    (click_count * 0.3) +
    (priority * 0.2) +
    (CASE
        WHEN updated_at >= NOW() - INTERVAL '7 days' THEN 10
        WHEN updated_at >= NOW() - INTERVAL '30 days' THEN 5
        ELSE 0
    END * 0.1)
)
WHERE is_active = true;

-- Add some test data for different languages (optional)
INSERT INTO autocomplete_entries (
    type, value, display_text, priority, search_count, click_count,
    synonyms, tags, score, language, metadata
) VALUES
('product', 'điện thoại iPhone', 'iPhone - Điện thoại thông minh Apple', 90, 80, 25,
 ARRAY['iphone', 'apple phone', 'smartphone'],
 ARRAY['electronics', 'apple', 'vietnamese'],
 85.0, 'vi', '{"price": 999, "currency": "USD", "category": "Electronics"}'),

('product', 'teléfono Samsung', 'Samsung Galaxy - Teléfono inteligente', 85, 70, 20,
 ARRAY['samsung', 'galaxy', 'smartphone'],
 ARRAY['electronics', 'samsung', 'spanish'],
 82.0, 'es', '{"price": 899, "currency": "USD", "category": "Electronics"}');

-- Verify the data was inserted correctly
SELECT
    type,
    COUNT(*) as count,
    AVG(score) as avg_score,
    COUNT(CASE WHEN is_trending THEN 1 END) as trending_count
FROM autocomplete_entries
WHERE is_active = true
GROUP BY type
ORDER BY count DESC;