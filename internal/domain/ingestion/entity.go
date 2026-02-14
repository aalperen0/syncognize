package ingestion

import (
	"time"

	"github.com/google/uuid"
)

// Ingestion Job
type IngestionStage string
type IngestionStatus string

const (
	IngestionStageUnspecified IngestionStage = "unspecified"
	IngestionStageClassifying IngestionStage = "classifying"
	IngestionStageChunking    IngestionStage = "chunking"
	IngestionStageDedup       IngestionStage = "dedup"
	IngestionStageEmbedding   IngestionStage = "embedding"
	IngestionStageStoring     IngestionStage = "storing"
)

const (
	IngestionStatusUnspecified IngestionStatus = "unspecified"
	IngestionStatusPending     IngestionStatus = "pending"
	IngestionStatusProcessing  IngestionStatus = "processing"
	IngestionStatusCompleted   IngestionStatus = "completed"
	IngestionStatusFailed      IngestionStatus = "failed"
	IngestionStatusRetrying    IngestionStatus = "retrying"
)

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
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type BufferEntry struct {
	Content   string
	Role      string
	Timestamp time.Time
}

type IngestionJob struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	MemoryID    uuid.UUID
	Status      IngestionStatus
	Stage       IngestionStage
	Error       *string
	Metadata    map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
}
