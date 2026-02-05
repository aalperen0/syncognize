CREATE SCHEMA IF NOT EXISTS ops;


-- =====================================
-- OPS SCHEMA: Operational/temporary
-- =====================================

-- Buffers (accumulate before flush to NATS)
CREATE TABLE ops.buffers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES init.users(id) ON DELETE CASCADE,
    source TEXT NOT NULL,
    entries JSONB NOT NULL DEFAULT '[]', -- Array of {content, role, timestamp}
    token_count INTEGER DEFAULT 0,
    word_count INTEGER DEFAULT 0,
    flush_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);



-- ================= INDEXES ================= 

CREATE INDEX idx_buffers_tenant_user ON ops.buffers(tenant_id, user_id);
CREATE INDEX idx_buffers_flush_at ON ops.buffers(flush_at);


-- Ingestion jobs (track async processing)
CREATE TABLE ops.ingestion_jobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    memory_id UUID REFERENCES init.memory(id) ON DELETE SET NULL,
    status TEXT NOT NULL CHECK(status IN ('pending', 'processing', 'completed', 'failed', 'retrying')),   -- pending, processing, completed, failed
    stage TEXT NOT NULL, -- chunking, embedding, dedup, storing, extracting
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


CREATE TRIGGER buffers_updated_at BEFORE UPDATE ON ops.buffers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER jobs_updated_at BEFORE UPDATE ON ops.ingestion_jobs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at();
