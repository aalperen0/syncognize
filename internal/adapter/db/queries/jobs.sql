-- name: CreateJob :exec
INSERT INTO ops.ingestion_jobs (
    id, tenant_id, memory_id, status, stage, metadata
) VALUES (
    $1, $2, $3, $4, $5, $6
);

-- name: GetJob :one
SELECT id, tenant_id, memory_id, status, stage,
       error, metadata, created_at, updated_at, completed_at
FROM ops.ingestion_jobs
WHERE id = $1;

-- name: UpdateJobStatus :exec
UPDATE ops.ingestion_jobs
SET status = $2, stage = $3
WHERE id = $1;

-- name: FailJob :exec
UPDATE ops.ingestion_jobs
SET status = 'failed', error = $2
WHERE id = $1;

-- name: CompleteJob :exec
UPDATE ops.ingestion_jobs
SET status = 'completed', completed_at = NOW(), metadata = jsonb_set(metadata, '{memory_ids}', $2::jsonb)
WHERE id = $1;
