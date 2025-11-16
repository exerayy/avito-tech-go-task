-- +goose Up
CREATE TYPE "pr_status" AS ENUM ('OPEN', 'MERGED');
-- +goose Down
DROP TYPE "pr_status";