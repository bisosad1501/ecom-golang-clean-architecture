-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create reviews table
CREATE TABLE reviews (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    product_id VARCHAR(36) NOT NULL,
    order_id VARCHAR(36) NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title VARCHAR(200) NOT NULL,
    comment TEXT,
    status ENUM('pending', 'approved', 'rejected') DEFAULT 'pending',
    is_verified BOOLEAN DEFAULT FALSE,
    helpful_count INT DEFAULT 0,
    not_helpful_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_reviews_user_id (user_id),
    INDEX idx_reviews_product_id (product_id),
    INDEX idx_reviews_order_id (order_id),
    INDEX idx_reviews_status (status),
    INDEX idx_reviews_rating (rating),
    INDEX idx_reviews_created_at (created_at),
    INDEX idx_reviews_is_verified (is_verified),
    
    UNIQUE KEY unique_user_product (user_id, product_id),
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE SET NULL
);

-- Create review_images table
CREATE TABLE review_images (
    id VARCHAR(36) PRIMARY KEY,
    review_id VARCHAR(36) NOT NULL,
    image_url VARCHAR(500) NOT NULL,
    image_alt VARCHAR(255),
    display_order INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    INDEX idx_review_images_review_id (review_id),
    INDEX idx_review_images_display_order (display_order),
    
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE
);

-- Create review_votes table
CREATE TABLE review_votes (
    id VARCHAR(36) PRIMARY KEY,
    review_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    vote_type ENUM('helpful', 'not_helpful') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_review_votes_review_id (review_id),
    INDEX idx_review_votes_user_id (user_id),
    INDEX idx_review_votes_vote_type (vote_type),
    
    UNIQUE KEY unique_user_review_vote (user_id, review_id),
    
    FOREIGN KEY (review_id) REFERENCES reviews(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create product_ratings table for aggregated rating data
CREATE TABLE product_ratings (
    id VARCHAR(36) PRIMARY KEY,
    product_id VARCHAR(36) NOT NULL UNIQUE,
    average_rating DECIMAL(3,2) DEFAULT 0.00,
    total_reviews INT DEFAULT 0,
    rating_1_count INT DEFAULT 0,
    rating_2_count INT DEFAULT 0,
    rating_3_count INT DEFAULT 0,
    rating_4_count INT DEFAULT 0,
    rating_5_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    
    INDEX idx_product_ratings_product_id (product_id),
    INDEX idx_product_ratings_average_rating (average_rating),
    INDEX idx_product_ratings_total_reviews (total_reviews),
    
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Create triggers to automatically update product ratings when reviews change
DELIMITER $$

CREATE TRIGGER update_product_rating_after_review_insert
AFTER INSERT ON reviews
FOR EACH ROW
BEGIN
    CALL update_product_rating_stats(NEW.product_id);
END$$

CREATE TRIGGER update_product_rating_after_review_update
AFTER UPDATE ON reviews
FOR EACH ROW
BEGIN
    IF OLD.rating != NEW.rating OR OLD.status != NEW.status THEN
        CALL update_product_rating_stats(NEW.product_id);
        IF OLD.product_id != NEW.product_id THEN
            CALL update_product_rating_stats(OLD.product_id);
        END IF;
    END IF;
END$$

CREATE TRIGGER update_product_rating_after_review_delete
AFTER DELETE ON reviews
FOR EACH ROW
BEGIN
    CALL update_product_rating_stats(OLD.product_id);
END$$

-- Create stored procedure to update product rating statistics
CREATE PROCEDURE update_product_rating_stats(IN p_product_id VARCHAR(36))
BEGIN
    DECLARE v_avg_rating DECIMAL(3,2) DEFAULT 0.00;
    DECLARE v_total_reviews INT DEFAULT 0;
    DECLARE v_rating_1 INT DEFAULT 0;
    DECLARE v_rating_2 INT DEFAULT 0;
    DECLARE v_rating_3 INT DEFAULT 0;
    DECLARE v_rating_4 INT DEFAULT 0;
    DECLARE v_rating_5 INT DEFAULT 0;
    
    -- Calculate statistics from approved reviews only
    SELECT 
        COALESCE(AVG(rating), 0.00),
        COUNT(*),
        SUM(CASE WHEN rating = 1 THEN 1 ELSE 0 END),
        SUM(CASE WHEN rating = 2 THEN 1 ELSE 0 END),
        SUM(CASE WHEN rating = 3 THEN 1 ELSE 0 END),
        SUM(CASE WHEN rating = 4 THEN 1 ELSE 0 END),
        SUM(CASE WHEN rating = 5 THEN 1 ELSE 0 END)
    INTO v_avg_rating, v_total_reviews, v_rating_1, v_rating_2, v_rating_3, v_rating_4, v_rating_5
    FROM reviews 
    WHERE product_id = p_product_id AND status = 'approved';
    
    -- Insert or update product rating record
    INSERT INTO product_ratings (
        id, product_id, average_rating, total_reviews,
        rating_1_count, rating_2_count, rating_3_count, rating_4_count, rating_5_count,
        created_at, updated_at
    ) VALUES (
        UUID(), p_product_id, v_avg_rating, v_total_reviews,
        v_rating_1, v_rating_2, v_rating_3, v_rating_4, v_rating_5,
        NOW(), NOW()
    )
    ON DUPLICATE KEY UPDATE
        average_rating = v_avg_rating,
        total_reviews = v_total_reviews,
        rating_1_count = v_rating_1,
        rating_2_count = v_rating_2,
        rating_3_count = v_rating_3,
        rating_4_count = v_rating_4,
        rating_5_count = v_rating_5,
        updated_at = NOW();
END$$

DELIMITER ;

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back

DROP TRIGGER IF EXISTS update_product_rating_after_review_delete;
DROP TRIGGER IF EXISTS update_product_rating_after_review_update;
DROP TRIGGER IF EXISTS update_product_rating_after_review_insert;
DROP PROCEDURE IF EXISTS update_product_rating_stats;

DROP TABLE IF EXISTS product_ratings;
DROP TABLE IF EXISTS review_votes;
DROP TABLE IF EXISTS review_images;
DROP TABLE IF EXISTS reviews;
