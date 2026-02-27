-- name: LinkEdge :exec
INSERT INTO graph.memory_edges (
    id, tenant_id, source_memory_id, target_memory_id,
    type, weight, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetNeighbors :many
SELECT id, tenant_id, source_memory_id, target_memory_id,
       type, weight, metadata, created_at
FROM graph.memory_edges
WHERE (source_memory_id = $1 OR target_memory_id = $1)
  AND tenant_id = $2;
