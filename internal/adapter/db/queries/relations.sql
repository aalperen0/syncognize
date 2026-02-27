-- name: StoreRelation :exec
INSERT INTO graph.entity_relations (
    id, tenant_id, source_entity_id, target_entity_id,
    predicate, memory_id, weight
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetRelationsByEntity :many
SELECT id, tenant_id, source_entity_id, target_entity_id,
       predicate, memory_id, weight, created_at
FROM graph.entity_relations
WHERE source_entity_id = $1 OR target_entity_id = $1;
