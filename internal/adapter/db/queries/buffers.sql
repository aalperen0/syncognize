-- name: GetBuffer :one
SELECT id, tenant_id, user_id, source, entries,
       token_count, flush_at, created_at, updated_at
FROM ops.buffers
WHERE tenant_id = $1 AND user_id = $2 AND source = $3;

-- name: CreateBuffer :exec
INSERT INTO ops.buffers (
    id, tenant_id, user_id, source, entries,
    token_count, flush_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: AppendBuffer :exec
UPDATE ops.buffers
SET entries = entries || $2::jsonb,
    token_count = token_count + $3
WHERE id = $1;

-- name: FlushBuffer :one
DELETE FROM ops.buffers
WHERE id = $1
RETURNING entries;

-- name: GetPendingFlush :many
SELECT id, tenant_id, user_id, source, entries,
       token_count, flush_at, created_at, updated_at
FROM ops.buffers
WHERE flush_at <= NOW();
