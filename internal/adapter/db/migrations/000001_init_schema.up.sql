CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;
CREATE EXTENSION IF NOT EXISTS citext;


CREATE SCHEMA IF NOT EXISTS init;
CREATE SCHEMA IF NOT EXISTS auth;

-- ==================================================================
-- TENANTS
-- ==================================================================

CREATE TABLE init.tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);



-- ==================================================================
-- USERS
-- ==================================================================


CREATE TABLE init.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    external_id TEXT NOT NULL,
    email CITEXT  NOT NULL,
    password_hash TEXT NOT NULL,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMPTZ,
    settings JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, email)

);


-- ================= INDEXES ================= 

CREATE INDEX idx_users_tenant_id ON init.users(tenant_id);
CREATE INDEX idx_users_email ON init.users(email);




-- ==================================================================
-- MEMORIES
-- ==================================================================

CREATE TABLE init.memories(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES init.users(id) ON DELETE CASCADE,
    source TEXT NOT NULL, --claude,chatgpt,gemini
    context_id TEXT NOT NULL, --conversation,buffer,import ID 
    type TEXT NOT NULL,
    content TEXT NOT NULL,
    content_hash TEXT NOT NULL,
    embedding vector(512),
    importance REAL DEFAULT 0.5,
    embedding_model TEXT DEFAULT 'voyage-4-lite',
    metadata JSONB default '{}',
    search_vector tsvector,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_accessed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ================= INDEXES ================= 



CREATE INDEX idx_memories_tenant_id ON init.memories(tenant_id);
CREATE INDEX idx_memories_user_id ON init.memories(user_id);
CREATE INDEX idx_memories_source ON init.memories(source);
CREATE INDEX idx_memories_context_id ON init.memories(context_id);
CREATE INDEX idx_memories_type ON init.memories(type);
CREATE INDEX idx_memories_content_hash ON init.memories(content_hash);

CREATE INDEX idx_memories__embedding ON init.memories
USING hnsw (embedding vector_cosine_ops)
WITH (m=16, ef_construction=64);

CREATE INDEX idx_memories_search_vector ON init.memories USING gin(search_vector);



-- ==================================================================
-- AUTH SCHEMA: API Keys
-- ==================================================================

CREATE TABLE auth.api_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    key_prefix TEXT NOT NULL,
    scopes TEXT[] DEFAULT '{"read", "write"}',
    rate_limit_rps INTEGER DEFAULT 10,
    is_active BOOLEAN DEFAULT true,
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by UUID REFERENCES init.users(id) ON DELETE SET NULL
);

-- ================= INDEXES ================= 
-- API Keys
CREATE INDEX idx_api_keys_tenant_id ON auth.api_keys(tenant_id);
CREATE INDEX idx_api_keys_key_prefix ON auth.api_keys(key_prefix);




-- ==================================================================
-- TRIGGERS & FUNCTIONS
-- ==================================================================

CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;    
$$ LANGUAGE plpgsql;


CREATE TRIGGER users_updated_at
BEFORE UPDATE ON init.users 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tenants_updated_at
BEFORE UPDATE ON init.tenants 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at();


-- Auto-generate search_vector for memories
CREATE OR REPLACE FUNCTION update_memory_search_vector()
RETURNS TRIGGER AS $$
BEGIN
     NEW.search_vector := to_tsvector('english', coalesce(NEW.content, ''));
     RETURN NEW;
END;
$$ LANGUAGE plpgsql;
 
CREATE TRIGGER memories_search_vector
BEFORE INSERT OR UPDATE OF content ON init.memories
FOR EACH ROW EXECUTE FUNCTION update_memory_search_vector();
