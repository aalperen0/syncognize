
-- ==================================================================
-- KNOWLEDGE GRAPH
-- ==================================================================


CREATE SCHEMA IF NOT EXISTS graph;


-- EXTRACTED ENTITIES                                                                
CREATE TABLE graph.extracted_entity(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    memory_id UUID NOT NULL REFERENCES init.memory(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    type TEXT NOT NULL, --person,org, loc, concept, event, product
    confidence REAL NOT NULL,
    aliases TEXT[] DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()
);


-- ================= INDEXES ================= 

      -- Extracted entities
CREATE INDEX idx_entities_tenant_id ON graph.extracted_entity(tenant_id);
CREATE INDEX idx_entities_memory_id ON graph.extracted_entity(memory_id);
CREATE INDEX idx_entities_name ON graph.extracted_entity(name);
CREATE INDEX idx_entities_type ON graph.extracted_entity(type);


-- ==================================================================
-- ENTITY-TO-ENTITY RELATIONS
-- ==================================================================



CREATE TABLE graph.entity_relations(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    source_entity_id UUID NOT NULL REFERENCES graph.extracted_entity(id) ON DELETE CASCADE,
    target_entity_id UUID NOT NULL REFERENCES graph.extracted_entity(id) ON DELETE CASCADE,
    predicate TEXT NOT NULL,    
    memory_id UUID NOT NULL REFERENCES init.memory(id) ON DELETE CASCADE,
    weight REAL DEFAULT 1.0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);


-- ================= INDEXES ================= 
-- Entity relations (graph traversal)
CREATE INDEX idx_entity_relations_source ON graph.entity_relations(source_entity_id);
CREATE INDEX idx_entity_relations_target ON graph.entity_relations(target_entity_id);
CREATE INDEX idx_entity_relations_predicate ON graph.entity_relations(predicate);
CREATE INDEX idx_entity_relations_tenant ON graph.entity_relations(tenant_id);



-- ==================================================================
-- MEMORY-TO-MEMORY EDGES(SEMANTIC CONNECTIONS)
-- ==================================================================

CREATE TABLE graph.memory_edges (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES init.tenants(id) ON DELETE CASCADE,
    source_memory_id UUID NOT NULL REFERENCES init.memory(id) ON DELETE CASCADE,
    target_memory_id UUID NOT NULL REFERENCES init.memory(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    weight REAL DEFAULT 1.0,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW()

);

-- ================= INDEXES ================= 

      
-- Memory edges (graph traversal)
CREATE INDEX idx_memory_edges_source ON graph.memory_edges(source_memory_id);
CREATE INDEX idx_memory_edges_target ON graph.memory_edges(target_memory_id);
CREATE INDEX idx_memory_edges_type ON graph.memory_edges(type);
CREATE INDEX idx_memory_edges_tenant ON graph.memory_edges(tenant_id);
