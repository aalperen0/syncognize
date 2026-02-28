package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/memory"
	"github.com/google/uuid"
)

type SqlcEdgeRepository struct {
	queries *sqlcgen.Queries
}

func NewEdgeRepository(queries *sqlcgen.Queries) *SqlcEdgeRepository {
	return &SqlcEdgeRepository{queries: queries}
}

func (r *SqlcEdgeRepository) Link(ctx context.Context, edge *memory.Edge) error {
	metadata := []byte("{}")
	if edge.Metadata != nil {
		var err error
		metadata, err = json.Marshal(edge.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal edge metadata: %w", err)
		}
	}

	return r.queries.LinkEdge(ctx, sqlcgen.LinkEdgeParams{
		ID:             edge.ID,
		TenantID:       edge.TenantID,
		SourceMemoryID: edge.SourceMemoryID,
		TargetMemoryID: edge.TargetMemoryID,
		Type:           string(edge.Type),
		Weight:         toFloat4(edge.Weight),
		Metadata:       metadata,
	})
}

func (r *SqlcEdgeRepository) Neighbors(ctx context.Context, memoryID uuid.UUID, tenantID uuid.UUID) ([]memory.Edge, error) {
	rows, err := r.queries.GetNeighbors(ctx, sqlcgen.GetNeighborsParams{
		SourceMemoryID: memoryID,
		TenantID:       tenantID,
	})
	if err != nil {
		return nil, err
	}

	edges := make([]memory.Edge, len(rows))
	for i, row := range rows {
		var metadata map[string]any
		if len(row.Metadata) > 0 {
			if err := json.Unmarshal(row.Metadata, &metadata); err != nil {
				return nil, fmt.Errorf("corrupt metadata for edge %s: %w", row.ID, err)
			}
		}

		edges[i] = memory.Edge{
			ID:             row.ID,
			TenantID:       row.TenantID,
			SourceMemoryID: row.SourceMemoryID,
			TargetMemoryID: row.TargetMemoryID,
			Type:           memory.EdgeType(row.Type),
			Weight:         toFloat32(row.Weight),
			Metadata:       metadata,
			CreatedAt:      toTime(row.CreatedAt),
		}
	}

	return edges, nil
}
