-- Drop triggers
DROP TRIGGER IF EXISTS jobs_updated_at ON ops.ingestion_jobs;

-- Drop tables (reverse dependency order)
DROP TABLE IF EXISTS ops.ingestion_jobs CASCADE;


DROP SCHEMA IF EXISTS ops;
