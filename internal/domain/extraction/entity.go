package extraction

import (
	"time"

	"github.com/google/uuid"
)

type EntityType string

const (
	EntityTypeUnspecified  EntityType = "unspecified"
	EntityTypePerson       EntityType = "person"
	EntityTypeOrganization EntityType = "organization"
	EntityTypeLocation     EntityType = "location"
	EntityTypeConcept      EntityType = "concept"
	EntityTypeEvent        EntityType = "event"
	EntityTypeProduct      EntityType = "product"
)

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
	TenantID       uuid.UUID
	SourceEntityID uuid.UUID
	TargetEntityID uuid.UUID
	Predicate      string
	MemoryID       uuid.UUID
	Weight         float32
	CreatedAt      time.Time
}
