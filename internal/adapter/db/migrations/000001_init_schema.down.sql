-- Drop triggers
DROP TRIGGER IF EXISTS memories_search_vector ON init.memory;
DROP TRIGGER IF EXISTS tenants_updated_at ON init.tenants;
DROP TRIGGER IF EXISTS users_updated_at ON init.users;

-- Drop tables (reverse dependency order)
DROP TABLE IF EXISTS init.memory CASCADE;
DROP TABLE IF EXISTS init.users CASCADE;
DROP TABLE IF EXISTS init.tenants CASCADE;
DROP TABLE IF EXISTS init.api_keys CASCADE;

-- Drop functions
DROP FUNCTION IF EXISTS update_memory_search_vector();
DROP FUNCTION IF EXISTS update_updated_at();

-- Drop extensions
DROP EXTENSION IF EXISTS unaccent;
DROP EXTENSION IF EXISTS pg_trgm;
DROP EXTENSION IF EXISTS vector;
DROP EXTENSION IF EXISTS citext;


DROP SCHEMA IF EXISTS init CASCADE;
