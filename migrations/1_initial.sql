-- +goose Up
CREATE TABLE IF NOT EXISTS "events" (
    "id" BIGINT PRIMARY KEY,
    "method" TEXT NOT NULL,
    "header" JSON NULL,
    "query_params" JSON NULL,
    "body" JSON NULL,
    "created_at" TIMESTAMPTZ DEFAULT 'NOW()',
    "room_id" UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "events";
