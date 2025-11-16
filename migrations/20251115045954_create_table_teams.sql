-- +goose Up
CREATE TABLE teams (
    name VARCHAR(36) PRIMARY KEY
);

-- +goose Down
DROP TABLE IF EXISTS teams;
