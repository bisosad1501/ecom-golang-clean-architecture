-- Remove product-tag associations
DELETE FROM product_tag_associations;

-- Remove brand associations from products
UPDATE products SET brand_id = NULL;
