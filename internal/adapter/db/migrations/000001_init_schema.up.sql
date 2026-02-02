CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS vector;
CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS unaccent;



-- ==================================================================
-- TENANTS
-- ==================================================================

CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(slug)
);


CREATE INDEX idx_tenants_slug ON tenants(slug);



-- ==================================================================
-- USERS
-- ==================================================================


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
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


CREATE INDEX idx_users_tenant_id ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);




-- ==================================================================
-- MEMORIES
-- ==================================================================

CREATE TABLE memories(
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
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

CREATE INDEX idx_memories_tenant_id ON memories(tenant_id);
CREATE INDEX idx_memories_user_id ON memories(user_id);
CREATE INDEX idx_memories_source ON memories(source);
CREATE INDEX idx_memories_context_id ON memories(context_id);
CREATE INDEX idx_memories_type ON memories(type);
CREATE INDEX idx_memories_content_hash ON memories(content_hash);

CREATE INDEX idx_memories__embedding ON memories
USING hnsw (embedding vector_cosine_ops)
WITH (m=16, ef_construction=64);

CREATE INDEX idx_memories_search_vector ON memories USING gin(search_vector);




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
BEFORE UPDATE ON users 
FOR EACH ROW 
EXECUTE FUNCTION update_updated_at();

CREATE TRIGGER tenants_updated_at
BEFORE UPDATE ON tenants 
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
BEFORE INSERT OR UPDATE OF content ON memories
FOR EACH ROW EXECUTE FUNCTION update_memory_search_vector();
