package memory

import (
	"time"

	"github.com/google/uuid"
)

type MemoryType string
type EdgeType string
type EntityType string

const (
	MemoryTypeUnspecified  MemoryType = "unspecified"
	MemoryTypeFact         MemoryType = "fact"
	MemoryTypePreference   MemoryType = "preference"
	MemoryTypeExperience   MemoryType = "experience"
	MemoryTypeSkill        MemoryType = "skill"
	MemoryTypeRelationship MemoryType = "relationship"
)

const (
	EdgeTypeUnspecified EdgeType = "unspecified"
	EdgeTypeRelatedTo   EdgeType = "related_to"
	EdgeTypeContradicts EdgeType = "contradicts"
	EdgeTypeUpdates     EdgeType = "updates"
)

const (
	EntityTypeUnspecified  EntityType = "unspecified"
	EntityTypePerson       EntityType = "person"
	EntityTypeOrganization EntityType = "organization"
	EntityTypeLocation     EntityType = "location"
	EntityTypeConcept      EntityType = "concept"
	EntityTypeEvent        EntityType = "event"
	EntityTypeProduct      EntityType = "product"
)

// Represents a stored semantic memory
// =======================================================
type Memory struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	UserID         uuid.UUID
	Source         string
	ContextId      string
	Type           MemoryType
	Content        string
	ContentHash    string
	Embedding      []float32
	Importance     float32
	EmbeddingModel string
	Metadata       map[string]interface{}
	CreatedAt      time.Time
	LastAccessedAt time.Time
}

// =======================================================
// Represents a relationship between memories
type Edge struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	SourceMemoryID uuid.UUID
	TargetMemoryID uuid.UUID
	Type           EdgeType
	Weight         float32
	Metadata       map[string]interface{}
}

// =======================================================
// Represents a named entity extracted from memory
type ExtractedEntity struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	MemoryID   uuid.UUID
	Name       string
	Type       EntityType
	Confidence float32
	Aliases    []string
	CreatedAt  time.Time
}

// =======================================================
// Represents a relationship between two entities
// "John works at X"
type EntityRelation struct {
	ID             uuid.UUID
	TenantId       uuid.UUID
	SourceEntityID uuid.UUID
	TargetEntityID uuid.UUID
	Predicate      string
	MemoryID       uuid.UUID
	Weight         float32
	CreatedAt      time.Time
}

// =======================================================
// Holds memories temporarily before batch processing
type Buffer struct {
	ID         uuid.UUID
	TenantID   uuid.UUID
	UserID     uuid.UUID
	Source     string
	Entries    []BufferEntry
	TokenCount int
	FlushAt    time.Time
}

type BufferEntry struct {
	Content   string
	Role      string
	Timestamp time.Time
}
