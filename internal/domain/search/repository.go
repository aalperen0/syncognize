package search

import (
	"context"

	"github.com/aalperen0/syncognize/internal/domain/memory"
)

type SearchRepository interface {
	SemanticSearch(ctx context.Context, embedding []float32, opts SearchFilters) ([]memory.ScoredMemory, error)
	KeywordSearch(ctx context.Context, query string, opts SearchFilters) ([]memory.ScoredMemory, error)
	HybridSearch(ctx context.Context, embedding []float32, query string, opts SearchFilters) ([]memory.ScoredMemory, error)
}
