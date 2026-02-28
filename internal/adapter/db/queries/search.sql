-- name: SemanticSearch :many
SELECT id, tenant_id, user_id, source, context_id,
       type, content, content_hash, content_type,
       language, scope, decay_rate, access_count,
       deleted_by, embedding, importance, embedding_model,
       metadata, created_at, last_accessed_at, deleted_at,
       (1 - (embedding <=> @embedding::vector))::real AS score
FROM init.memory
WHERE tenant_id = @tenant_id
  AND deleted_at IS NULL
  AND (CASE WHEN @filter_user::bool THEN user_id = @user_id ELSE true END)
  AND (CASE WHEN @filter_types::bool THEN type = ANY(@memory_types::text[]) ELSE true END)
  AND (CASE WHEN @filter_sources::bool THEN source = ANY(@sources::text[]) ELSE true END)
  AND (CASE WHEN @filter_content_types::bool THEN content_type = ANY(@content_types::text[]) ELSE true END)
  AND (CASE WHEN @filter_start_date::bool THEN created_at >= @start_date ELSE true END)
  AND (CASE WHEN @filter_end_date::bool THEN created_at <= @end_date ELSE true END)
ORDER BY embedding <=> @embedding::vector
LIMIT @result_limit;

-- name: KeywordSearch :many
SELECT id, tenant_id, user_id, source, context_id,
       type, content, content_hash, content_type,
       language, scope, decay_rate, access_count,
       deleted_by, embedding, importance, embedding_model,
       metadata, created_at, last_accessed_at, deleted_at,
       ts_rank(search_vector, plainto_tsquery('english', @query)) AS score
FROM init.memory
WHERE tenant_id = @tenant_id
  AND deleted_at IS NULL
  AND search_vector @@ plainto_tsquery('english', @query)
  AND (CASE WHEN @filter_user::bool THEN user_id = @user_id ELSE true END)
  AND (CASE WHEN @filter_types::bool THEN type = ANY(@memory_types::text[]) ELSE true END)
  AND (CASE WHEN @filter_sources::bool THEN source = ANY(@sources::text[]) ELSE true END)
  AND (CASE WHEN @filter_content_types::bool THEN content_type = ANY(@content_types::text[]) ELSE true END)
  AND (CASE WHEN @filter_start_date::bool THEN created_at >= @start_date ELSE true END)
  AND (CASE WHEN @filter_end_date::bool THEN created_at <= @end_date ELSE true END)
ORDER BY score DESC
LIMIT @result_limit;
