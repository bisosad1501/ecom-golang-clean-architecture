-- Update products with brand associations
UPDATE products SET brand_id = (SELECT id FROM brands WHERE name = 'Apple' LIMIT 1) 
WHERE name IN ('iPhone 15 Pro', 'MacBook Air M3');

UPDATE products SET brand_id = (SELECT id FROM brands WHERE name = 'Samsung' LIMIT 1) 
WHERE name = 'Samsung Galaxy S24';

UPDATE products SET brand_id = (SELECT id FROM brands WHERE name = 'Sony' LIMIT 1) 
WHERE name = 'Sony WH-1000XM5';

UPDATE products SET brand_id = (SELECT id FROM brands WHERE name = 'Nike' LIMIT 1) 
WHERE name IN ('Nike Air Max 270', 'Nike Dri-FIT T-Shirt');

UPDATE products SET brand_id = (SELECT id FROM brands WHERE name = 'Nike' LIMIT 1)
WHERE name = 'Levi''s 501 Jeans';

-- Insert product-tag associations
-- Electronics products with Featured tag
INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'iPhone 15 Pro' AND t.name = 'Featured';

INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'Samsung Galaxy S24' AND t.name = 'Popular';

INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'MacBook Air M3' AND t.name = 'Featured';

INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'Sony WH-1000XM5' AND t.name = 'Popular';

-- Clothing products with New tag
INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'Nike Air Max 270' AND t.name = 'New';

INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'Nike Dri-FIT T-Shirt' AND t.name = 'Sale';

INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name = 'Levi''s 501 Jeans' AND t.name = 'Popular';

-- Books with Limited Edition tags
INSERT INTO product_tag_associations (product_id, product_tag_id)
SELECT p.id, t.id
FROM products p, tags t
WHERE p.name IN ('Clean Code', 'The Pragmatic Programmer', 'Design Patterns') AND t.name = 'Limited Edition';
