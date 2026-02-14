package memory

import (
	"context"

	"github.com/google/uuid"
)

type MemoryRepository interface {
	GetMemory(ctx context.Context, memoryID uuid.UUID, tenantID uuid.UUID) (*Memory, error)
	GetByContentHash(ctx context.Context, hash string, tenantID uuid.UUID) (*Memory, error)
	Store(ctx context.Context, m *Memory) error
	Delete(ctx context.Context, memoryID uuid.UUID, tenantID uuid.UUID) error
	Update(ctx context.Context, m *Memory) error
}

type EdgeRepository interface {
	Link(ctx context.Context, edge *Edge) error
	Neighbors(ctx context.Context, memoryID uuid.UUID) ([]Edge, error)
}
