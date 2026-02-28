package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/ingestion"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SqlcBufferRepository struct {
	queries *sqlcgen.Queries
}

func NewBufferRepository(queries *sqlcgen.Queries) *SqlcBufferRepository {
	return &SqlcBufferRepository{queries: queries}
}

func (r *SqlcBufferRepository) Get(ctx context.Context, tenantID, userID uuid.UUID, source string) (*ingestion.Buffer, error) {
	row, err := r.queries.GetBuffer(ctx, sqlcgen.GetBufferParams{
		TenantID: tenantID,
		UserID:   userID,
		Source:   source,
	})
	if err != nil {
		return nil, err
	}
	return bufferToDomain(row)
}

func (r *SqlcBufferRepository) Append(ctx context.Context, tenantID, userID uuid.UUID, source string, entry ingestion.BufferEntry) error {
	buf, err := r.queries.GetBuffer(ctx, sqlcgen.GetBufferParams{
		TenantID: tenantID,
		UserID:   userID,
		Source:   source,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return r.createWithEntry(ctx, tenantID, userID, source, entry)
		}
		return err
	}

	entryJSON, err := json.Marshal([]ingestion.BufferEntry{entry})
	if err != nil {
		return fmt.Errorf("failed to marshal buffer entry: %w", err)
	}

	return r.queries.AppendBuffer(ctx, sqlcgen.AppendBufferParams{
		ID:         buf.ID,
		Column2:    entryJSON,
		TokenCount: toInt4(int32(len(entry.Content) / 4)), // rough token estimate
	})
}

func (r *SqlcBufferRepository) Flush(ctx context.Context, bufferID uuid.UUID) ([]ingestion.BufferEntry, error) {
	entriesJSON, err := r.queries.FlushBuffer(ctx, bufferID)
	if err != nil {
		return nil, err
	}

	var entries []ingestion.BufferEntry
	if len(entriesJSON) > 0 {
		if err := json.Unmarshal(entriesJSON, &entries); err != nil {
			return nil, fmt.Errorf("corrupt buffer entries for buffer %s: %w", bufferID, err)
		}
	}

	return entries, nil
}

func (r *SqlcBufferRepository) GetPendingFlush(ctx context.Context) ([]ingestion.Buffer, error) {
	rows, err := r.queries.GetPendingFlush(ctx)
	if err != nil {
		return nil, err
	}

	buffers := make([]ingestion.Buffer, len(rows))
	for i, row := range rows {
		buf, err := bufferToDomain(row)
		if err != nil {
			return nil, err
		}
		buffers[i] = *buf
	}

	return buffers, nil
}

func (r *SqlcBufferRepository) createWithEntry(ctx context.Context, tenantID, userID uuid.UUID, source string, entry ingestion.BufferEntry) error {
	entriesJSON, err := json.Marshal([]ingestion.BufferEntry{entry})
	if err != nil {
		return fmt.Errorf("failed to marshal buffer entry: %w", err)
	}

	return r.queries.CreateBuffer(ctx, sqlcgen.CreateBufferParams{
		ID:         uuid.New(),
		TenantID:   tenantID,
		UserID:     userID,
		Source:     source,
		Entries:    entriesJSON,
		TokenCount: toInt4(int32(len(entry.Content) / 4)),
		FlushAt:    toTimestampz(entry.Timestamp.Add(5 * 60 * 1e9)), // flush_at = now + 5min
	})
}

func bufferToDomain(row sqlcgen.OpsBuffer) (*ingestion.Buffer, error) {
	var entries []ingestion.BufferEntry
	if len(row.Entries) > 0 {
		if err := json.Unmarshal(row.Entries, &entries); err != nil {
			return nil, fmt.Errorf("corrupt buffer entries for buffer %s: %w", row.ID, err)
		}
	}

	return &ingestion.Buffer{
		ID:         row.ID,
		TenantID:   row.TenantID,
		UserID:     row.UserID,
		Source:     row.Source,
		Entries:    entries,
		TokenCount: int(toInt32(row.TokenCount)),
		FlushAt:    toTime(row.FlushAt),
		CreatedAt:  toTime(row.CreatedAt),
		UpdatedAt:  toTime(row.UpdatedAt),
	}, nil
}
