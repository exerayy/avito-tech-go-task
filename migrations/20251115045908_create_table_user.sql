-- +goose Up
CREATE TABLE users (
    id        VARCHAR(36) PRIMARY KEY,
    name      VARCHAR(255) NOT NULL,
    team_name   VARCHAR(36)  NOT NULL,
    is_active BOOLEAN DEFAULT TRUE
);

-- +goose Down
DROP TABLE IF EXISTS users;
