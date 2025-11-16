-- +goose Up
CREATE TABLE user_review_stats (
    user_id VARCHAR(36) PRIMARY KEY,
    total_reviews INT DEFAULT 0,
    active_reviews INT DEFAULT 0,
    merged_reviews INT DEFAULT 0,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_review_stats;
