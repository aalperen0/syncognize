-- name: StoreMemory :exec
INSERT INTO init.memory (
    id, tenant_id, user_id, source, context_id,
    type, content, content_hash, content_type,
    language, scope, decay_rate, embedding,
    importance, embedding_model, metadata
) VALUES (
    $1, $2, $3, $4, $5,
    $6, $7, $8, $9,
    $10, $11, $12, $13,
    $14, $15, $16
);

-- name: GetMemory :one
SELECT id, tenant_id, user_id, source, context_id,
       type, content, content_hash, content_type,
       language, scope, decay_rate, access_count,
       deleted_by, embedding, importance, embedding_model,
       metadata, created_at, last_accessed_at, deleted_at
FROM init.memory
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL;

-- name: GetByContentHash :one
SELECT id, tenant_id, user_id, source, context_id,
       type, content, content_hash, content_type,
       language, scope, decay_rate, access_count,
       deleted_by, embedding, importance, embedding_model,
       metadata, created_at, last_accessed_at, deleted_at
FROM init.memory
WHERE content_hash = $1 AND tenant_id = $2 AND deleted_at IS NULL;

-- name: UpdateMemory :exec
UPDATE init.memory SET
    content = $3,
    content_hash = $4,
    content_type = $5,
    language = $6,
    scope = $7,
    importance = $8,
    access_count = $9,
    metadata = $10,
    last_accessed_at = NOW()
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL;

-- name: DeleteMemory :exec
UPDATE init.memory
SET deleted_at = NOW(), deleted_by = $3
WHERE id = $1 AND tenant_id = $2 AND deleted_at IS NULL;
