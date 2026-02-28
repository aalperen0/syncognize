CREATE SCHEMA IF NOT EXISTS ops;


-- =====================================
-- OPS SCHEMA: Operational/temporary
-- =====================================

-- Ingestion jobs (track async processing)
CREATE TABLE ops.ingestion_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    memory_id UUID REFERENCES init.memory(id) ON DELETE SET NULL,
    status TEXT NOT NULL CHECK(status IN ('unspecified','pending', 'processing', 'completed', 'failed', 'retrying')),   -- pending, processing, completed, failed
    stage TEXT NOT NULL CHECK(stage IN ('unspecified', 'classifying', 'chunking', 'dedup', 'embedding', 'storing')), -- chunking, embedding, dedup, storing, extracting
    error TEXT,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    completed_at TIMESTAMPTZ
);


-- ================= INDEXES =================
-- Ingestion jobs
CREATE INDEX idx_jobs_status ON ops.ingestion_jobs(status);
CREATE INDEX idx_jobs_memory_id ON ops.ingestion_jobs(memory_id);
CREATE INDEX idx_jobs_created_at ON ops.ingestion_jobs(created_at DESC);

CREATE TRIGGER jobs_updated_at BEFORE UPDATE ON ops.ingestion_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
