-- +goose Up
CREATE TABLE IF NOT EXISTS "events" (
    "id" BIGINT PRIMARY KEY,
    "method" TEXT NOT NULL,
    "header" JSON NOT NULL DEFAULT '{}',
    "query_params" JSON NOT NULL DEFAULT '{}',
    "body" JSON NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMPTZ DEFAULT 'NOW()',
    "room_id" UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS "events";
