-- +goose Up
CREATE TABLE pull_requests (
   id VARCHAR(36) PRIMARY KEY,
   name VARCHAR(255) NOT NULL,
   author_id VARCHAR(36) NOT NULL,
   status pr_status NOT NULL,
   reviewers_ids TEXT[] NOT NULL DEFAULT '{}',
   merged_at TIMESTAMP WITH TIME ZONE
);

-- +goose Down
DROP TABLE IF EXISTS pull_requests;
