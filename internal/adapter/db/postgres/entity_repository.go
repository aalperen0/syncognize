package postgres

import (
	"context"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/extraction"
	"github.com/google/uuid"
)

type SqlcEntityRepository struct {
	queries *sqlcgen.Queries
}

func NewEntityRepository(queries *sqlcgen.Queries) *SqlcEntityRepository {
	return &SqlcEntityRepository{queries: queries}
}

func (r *SqlcEntityRepository) Store(ctx context.Context, e *extraction.ExtractedEntity) error {
	return r.queries.StoreEntity(ctx, sqlcgen.StoreEntityParams{
		ID:         e.ID,
		TenantID:   e.TenantID,
		MemoryID:   e.MemoryID,
		Name:       e.Name,
		Type:       string(e.Type),
		Confidence: e.Confidence,
		Aliases:    e.Aliases,
	})
}

func (r *SqlcEntityRepository) GetByMemory(ctx context.Context, memoryID uuid.UUID) ([]extraction.ExtractedEntity, error) {
	rows, err := r.queries.GetEntitiesByMemory(ctx, memoryID)
	if err != nil {
		return nil, err
	}
	return entitiesToDomain(rows), nil
}

func (r *SqlcEntityRepository) FindByName(ctx context.Context, name string) ([]extraction.ExtractedEntity, error) {
	rows, err := r.queries.FindEntitiesByName(ctx, name)
	if err != nil {
		return nil, err
	}
	return entitiesToDomain(rows), nil
}

func entitiesToDomain(rows []sqlcgen.GraphExtractedEntity) []extraction.ExtractedEntity {
	entities := make([]extraction.ExtractedEntity, len(rows))
	for i, row := range rows {
		entities[i] = extraction.ExtractedEntity{
			ID:         row.ID,
			TenantID:   row.TenantID,
			MemoryID:   row.MemoryID,
			Name:       row.Name,
			Type:       extraction.EntityType(row.Type),
			Confidence: row.Confidence,
			Aliases:    row.Aliases,
			CreatedAt:  toTime(row.CreatedAt),
		}
	}
	return entities
}
