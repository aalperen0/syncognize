package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/ingestion"
	"github.com/google/uuid"
)

type SqlcIngestionRepository struct {
	queries *sqlcgen.Queries
}

func NewIngestionRepository(queries *sqlcgen.Queries) *SqlcIngestionRepository {
	return &SqlcIngestionRepository{queries: queries}
}

func (r *SqlcIngestionRepository) CreateJob(ctx context.Context, job *ingestion.IngestionJob) error {
	metadata := []byte("{}")
	if job.Metadata != nil {
		var err error
		metadata, err = json.Marshal(job.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal job metadata: %w", err)
		}
	}

	return r.queries.CreateJob(ctx, sqlcgen.CreateJobParams{
		ID:       job.ID,
		TenantID: job.TenantID,
		MemoryID: toPgUUID(&job.MemoryID),
		Status:   string(job.Status),
		Stage:    string(job.Stage),
		Metadata: metadata,
	})
}

func (r *SqlcIngestionRepository) GetJob(ctx context.Context, jobID uuid.UUID) (*ingestion.IngestionJob, error) {
	row, err := r.queries.GetJob(ctx, jobID)
	if err != nil {
		return nil, err
	}
	return jobToDomain(row)
}

func (r *SqlcIngestionRepository) UpdateStatus(ctx context.Context, jobID uuid.UUID, status ingestion.IngestionStatus, stage ingestion.IngestionStage) error {
	return r.queries.UpdateJobStatus(ctx, sqlcgen.UpdateJobStatusParams{
		ID:     jobID,
		Status: string(status),
		Stage:  string(stage),
	})
}

func (r *SqlcIngestionRepository) FailJob(ctx context.Context, jobID uuid.UUID, errMsg string) error {
	return r.queries.FailJob(ctx, sqlcgen.FailJobParams{
		ID:    jobID,
		Error: toText(errMsg),
	})
}

func (r *SqlcIngestionRepository) CompleteJob(ctx context.Context, jobID uuid.UUID, memoryIDs []uuid.UUID) error {
	idsJSON, err := json.Marshal(memoryIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal memory IDs: %w", err)
	}

	return r.queries.CompleteJob(ctx, sqlcgen.CompleteJobParams{
		ID:      jobID,
		Column2: idsJSON,
	})
}

func jobToDomain(row sqlcgen.OpsIngestionJob) (*ingestion.IngestionJob, error) {
	var metadata map[string]any
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("corrupt metadata for job %s: %w", row.ID, err)
		}
	}

	var errStr *string
	if row.Error.Valid {
		errStr = &row.Error.String
	}

	memoryID := uuid.Nil
	if row.MemoryID.Valid {
		memoryID = uuid.UUID(row.MemoryID.Bytes)
	}

	return &ingestion.IngestionJob{
		ID:          row.ID,
		TenantID:    row.TenantID,
		MemoryID:    memoryID,
		Status:      ingestion.IngestionStatus(row.Status),
		Stage:       ingestion.IngestionStage(row.Stage),
		Error:       errStr,
		Metadata:    metadata,
		CreatedAt:   toTime(row.CreatedAt),
		UpdatedAt:   toTime(row.UpdatedAt),
		CompletedAt: toTime(row.CompletedAt),
	}, nil
}
