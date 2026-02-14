package ingestion

import (
	"context"

	"github.com/google/uuid"
)

type BufferRepository interface {
	Get(ctx context.Context, tenantID, userID uuid.UUID, source string) (*Buffer, error)
	Append(ctx context.Context, tenantID, userID uuid.UUID, source string, entry BufferEntry) error
	Flush(ctx context.Context, bufferID uuid.UUID) ([]BufferEntry, error)
	GetPendingFlush(ctx context.Context) ([]Buffer, error)
}

type IngestionRepository interface {
	CreateJob(ctx context.Context, job *IngestionJob) error
	GetJob(ctx context.Context, jobID uuid.UUID) (*IngestionJob, error)
	UpdateStatus(ctx context.Context, jobID uuid.UUID, status IngestionStatus, stage IngestionStage) error
	FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error
	CompleteJob(ctx context.Context, jobID uuid.UUID, memoryIDs []uuid.UUID) error
}
