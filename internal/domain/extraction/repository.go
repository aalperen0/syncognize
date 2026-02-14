package extraction

import (
	"context"

	"github.com/google/uuid"
)

type ExtractedEntityRepository interface {
	GetByMemory(ctx context.Context, memoryID uuid.UUID) ([]ExtractedEntity, error)
	FindByName(ctx context.Context, name string) ([]ExtractedEntity, error)
	Store(ctx context.Context, e *ExtractedEntity) error
}

type EntityRelationRepository interface {
	GetByEntity(ctx context.Context, entityID uuid.UUID) ([]EntityRelation, error)
	Store(ctx context.Context, er *EntityRelation) error
}
