package memory

import (
	"time"

	"github.com/google/uuid"
)

type MemoryType string
type EdgeType string
type ContentType string

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
	ContentTypeUnspecified ContentType = "unspecified"
	ContentTypeText        ContentType = "text"
	ContentTypeCode        ContentType = "code"
	ContentTypeMixed       ContentType = "mixed"
)

// Represents a stored semantic memory
// =======================================================
type Memory struct {
	ID             uuid.UUID
	TenantID       uuid.UUID
	UserID         uuid.UUID
	Source         string
	ContextID      string
	Type           MemoryType
	Content        string
	ContentHash    string
	ContentType    ContentType // text, code, mixed
	Language       string      // for code chunks: go, python, etc
	Scope          string
	DecayRate      float32
	AccessCount    int32
	DeletedBy      *uuid.UUID
	Embedding      []float32
	Importance     float32
	EmbeddingModel string
	Metadata       map[string]any
	CreatedAt      time.Time
	LastAccessedAt time.Time
	DeletedAt      *time.Time
}

type ScoredMemory struct {
	Memory      Memory
	Score       float32
	MatchSource string
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
	Metadata       map[string]any
	CreatedAt      time.Time
}
