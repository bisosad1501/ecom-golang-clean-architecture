-- Drop triggers
DROP TRIGGER IF EXISTS update_frequently_bought_together_updated_at ON frequently_bought_together;
DROP TRIGGER IF EXISTS update_product_similarities_updated_at ON product_similarities;
DROP TRIGGER IF EXISTS update_product_recommendations_updated_at ON product_recommendations;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order
DROP TABLE IF EXISTS trending_products;
DROP TABLE IF EXISTS frequently_bought_together;
DROP TABLE IF EXISTS product_similarities;
DROP TABLE IF EXISTS user_product_interactions;
DROP TABLE IF EXISTS product_recommendations;
