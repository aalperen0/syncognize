package postgres

import (
	"context"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/extraction"
	"github.com/google/uuid"
)

type SqlcRelationRepository struct {
	queries *sqlcgen.Queries
}

func NewRelationRepository(queries *sqlcgen.Queries) *SqlcRelationRepository {
	return &SqlcRelationRepository{queries: queries}
}

func (r *SqlcRelationRepository) Store(ctx context.Context, er *extraction.EntityRelation) error {
	return r.queries.StoreRelation(ctx, sqlcgen.StoreRelationParams{
		ID:             er.ID,
		TenantID:       er.TenantID,
		SourceEntityID: er.SourceEntityID,
		TargetEntityID: er.TargetEntityID,
		Predicate:      er.Predicate,
		MemoryID:       er.MemoryID,
		Weight:         toFloat4(er.Weight),
	})
}

func (r *SqlcRelationRepository) GetByEntity(ctx context.Context, entityID uuid.UUID) ([]extraction.EntityRelation, error) {
	rows, err := r.queries.GetRelationsByEntity(ctx, entityID)
	if err != nil {
		return nil, err
	}

	relations := make([]extraction.EntityRelation, len(rows))
	for i, row := range rows {
		relations[i] = extraction.EntityRelation{
			ID:             row.ID,
			TenantID:       row.TenantID,
			SourceEntityID: row.SourceEntityID,
			TargetEntityID: row.TargetEntityID,
			Predicate:      row.Predicate,
			MemoryID:       row.MemoryID,
			Weight:         toFloat32(row.Weight),
			CreatedAt:      toTime(row.CreatedAt),
		}
	}

	return relations, nil
}
