-- name: SaveEvent :one
INSERT INTO events (id, method, header, query_params, body, room_id)
VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (id) DO
UPDATE SET method = EXCLUDED.method,
    header = EXCLUDED.header,
    query_params = EXCLUDED.query_params,
    body = EXCLUDED.body,
    room_id = EXCLUDED.room_id
RETURNING created_at;

-- name: ListEvents :many
SELECT * FROM events WHERE room_id = $1 ORDER BY created_at DESC OFFSET $2 LIMIT $3;
