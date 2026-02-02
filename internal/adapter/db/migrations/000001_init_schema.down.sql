-- Drop triggers
DROP TRIGGER IF EXISTS memories_search_vector ON memories;
DROP TRIGGER IF EXISTS tenants_updated_at ON tenants;
DROP TRIGGER IF EXISTS users_updated_at ON users;

-- Drop tables (reverse dependency order)
DROP TABLE IF EXISTS memories CASCADE;
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS tenants CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_memory_search_vector();
DROP FUNCTION IF EXISTS update_updated_at();

-- Drop extensions
DROP EXTENSION IF EXISTS unaccent;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS "uuid-ossp";
