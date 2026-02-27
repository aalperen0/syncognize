-- name: StoreEntity :exec
INSERT INTO graph.extracted_entity (
    id, tenant_id, memory_id, name, type, confidence, aliases
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetEntitiesByMemory :many
SELECT id, tenant_id, memory_id, name, type,
       confidence, aliases, created_at
FROM graph.extracted_entity
WHERE memory_id = $1;

-- name: FindEntitiesByName :many
SELECT id, tenant_id, memory_id, name, type,
       confidence, aliases, created_at
FROM graph.extracted_entity
WHERE name ILIKE $1;
