package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/memory"
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

type SqlcMemoryRepository struct {
	queries *sqlcgen.Queries
}

func NewMemoryRepository(queries *sqlcgen.Queries) *SqlcMemoryRepository {
	return &SqlcMemoryRepository{queries: queries}
}

func (r *SqlcMemoryRepository) GetMemory(ctx context.Context, memoryID uuid.UUID, tenantID uuid.UUID) (*memory.Memory, error) {
	row, err := r.queries.GetMemory(ctx, sqlcgen.GetMemoryParams{
		ID:       memoryID,
		TenantID: tenantID,
	})
	if err != nil {
		return nil, err
	}
	return toDomain(row)
}

func (r *SqlcMemoryRepository) Store(ctx context.Context, m *memory.Memory) error {
	metadata := []byte("{}")
	if m.Metadata != nil {
		var err error
		metadata, err = json.Marshal(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	return r.queries.StoreMemory(ctx, sqlcgen.StoreMemoryParams{
		ID:             m.ID,
		TenantID:       m.TenantID,
		UserID:         m.UserID,
		Source:         m.Source,
		ContextID:      m.ContextID,
		Type:           string(m.Type),
		Content:        m.Content,
		ContentHash:    m.ContentHash,
		ContentType:    string(m.ContentType),
		Language:       toText(m.Language),
		Scope:          toText(m.Scope),
		DecayRate:      toFloat4(m.DecayRate),
		Embedding:      pgvector.NewVector(m.Embedding),
		Importance:     toFloat4(m.Importance),
		EmbeddingModel: toText(m.EmbeddingModel),
		Metadata:       metadata,
	})
}

func (r *SqlcMemoryRepository) GetByContentHash(ctx context.Context, hash string, tenantID uuid.UUID) (*memory.Memory, error) {
	row, err := r.queries.GetByContentHash(ctx, sqlcgen.GetByContentHashParams{
		ContentHash: hash,
		TenantID:    tenantID,
	})
	if err != nil {
		return nil, err
	}
	return toDomainFromHash(row)
}

func (r *SqlcMemoryRepository) Update(ctx context.Context, m *memory.Memory) error {
	metadata := []byte("{}")
	if m.Metadata != nil {
		var err error
		metadata, err = json.Marshal(m.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
	}

	return r.queries.UpdateMemory(ctx, sqlcgen.UpdateMemoryParams{
		ID:          m.ID,
		TenantID:    m.TenantID,
		Content:     m.Content,
		ContentHash: m.ContentHash,
		ContentType: string(m.ContentType),
		Language:    toText(m.Language),
		Scope:       toText(m.Scope),
		Importance:  toFloat4(m.Importance),
		AccessCount: toInt4(m.AccessCount),
		Metadata:    metadata,
	})
}

func (r *SqlcMemoryRepository) Delete(ctx context.Context, memoryID uuid.UUID, tenantID uuid.UUID) error {
	return r.queries.DeleteMemory(ctx, sqlcgen.DeleteMemoryParams{
		ID:        memoryID,
		TenantID:  tenantID,
		DeletedBy: toPgUUID(nil),
	})
}

func toDomain(row sqlcgen.GetMemoryRow) (*memory.Memory, error) {
	var metadata map[string]any
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("corrupt metadata for memory %s: %w", row.ID, err)
		}
	}

	return &memory.Memory{
		ID:             row.ID,
		TenantID:       row.TenantID,
		UserID:         row.UserID,
		Source:         row.Source,
		ContextID:      row.ContextID,
		Type:           memory.MemoryType(row.Type),
		Content:        row.Content,
		ContentHash:    row.ContentHash,
		ContentType:    memory.ContentType(row.ContentType),
		Language:       toString(row.Language),
		Scope:          toString(row.Scope),
		DecayRate:      toFloat32(row.DecayRate),
		AccessCount:    toInt32(row.AccessCount),
		DeletedBy:      toUUIDPtr(row.DeletedBy),
		Embedding:      row.Embedding.Slice(),
		Importance:     toFloat32(row.Importance),
		EmbeddingModel: toString(row.EmbeddingModel),
		Metadata:       metadata,
		CreatedAt:      toTime(row.CreatedAt),
		LastAccessedAt: toTime(row.LastAccessedAt),
		DeletedAt:      toTimePtr(row.DeletedAt),
	}, nil
}

func toDomainFromHash(row sqlcgen.GetByContentHashRow) (*memory.Memory, error) {
	var metadata map[string]any
	if len(row.Metadata) > 0 {
		if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
			return nil, fmt.Errorf("corrupt metadata for memory %s: %w", row.ID, err)
		}
	}

	return &memory.Memory{
		ID:             row.ID,
		TenantID:       row.TenantID,
		UserID:         row.UserID,
		Source:         row.Source,
		ContextID:      row.ContextID,
		Type:           memory.MemoryType(row.Type),
		Content:        row.Content,
		ContentHash:    row.ContentHash,
		ContentType:    memory.ContentType(row.ContentType),
		Language:       toString(row.Language),
		Scope:          toString(row.Scope),
		DecayRate:      toFloat32(row.DecayRate),
		AccessCount:    toInt32(row.AccessCount),
		DeletedBy:      toUUIDPtr(row.DeletedBy),
		Embedding:      row.Embedding.Slice(),
		Importance:     toFloat32(row.Importance),
		EmbeddingModel: toString(row.EmbeddingModel),
		Metadata:       metadata,
		CreatedAt:      toTime(row.CreatedAt),
		LastAccessedAt: toTime(row.LastAccessedAt),
		DeletedAt:      toTimePtr(row.DeletedAt),
	}, nil
}
