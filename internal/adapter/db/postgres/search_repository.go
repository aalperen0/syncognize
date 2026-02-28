package postgres

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/aalperen0/syncognize/internal/adapter/db/sqlcgen"
	"github.com/aalperen0/syncognize/internal/domain/memory"
	"github.com/aalperen0/syncognize/internal/domain/search"
	"github.com/google/uuid"
	pgvector "github.com/pgvector/pgvector-go"
)

const (
	defaultLimit = 20
	rrfK         = 60
)

type SqlcSearchRepository struct {
	queries *sqlcgen.Queries
}

func NewSearchRepository(queries *sqlcgen.Queries) *SqlcSearchRepository {
	return &SqlcSearchRepository{queries: queries}
}

func (r *SqlcSearchRepository) SemanticSearch(ctx context.Context, embedding []float32, opts search.SearchFilters) ([]memory.ScoredMemory, error) {
	params := buildSemanticParams(embedding, opts)
	rows, err := r.queries.SemanticSearch(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]memory.ScoredMemory, len(rows))
	for i, row := range rows {
		m, err := semanticRowToDomain(row)
		if err != nil {
			return nil, err
		}
		results[i] = memory.ScoredMemory{
			Memory:      *m,
			Score:       float32(row.Score),
			MatchSource: "semantic",
		}
	}
	return results, nil
}

func (r *SqlcSearchRepository) KeywordSearch(ctx context.Context, query string, opts search.SearchFilters) ([]memory.ScoredMemory, error) {
	params := buildKeywordParams(query, opts)
	rows, err := r.queries.KeywordSearch(ctx, params)
	if err != nil {
		return nil, err
	}

	results := make([]memory.ScoredMemory, len(rows))
	for i, row := range rows {
		m, err := keywordRowToDomain(row)
		if err != nil {
			return nil, err
		}
		results[i] = memory.ScoredMemory{
			Memory:      *m,
			Score:       row.Score,
			MatchSource: "keyword",
		}
	}
	return results, nil
}

func (r *SqlcSearchRepository) HybridSearch(ctx context.Context, embedding []float32, query string, opts search.SearchFilters) ([]memory.ScoredMemory, error) {
	semanticResults, err := r.SemanticSearch(ctx, embedding, opts)
	if err != nil {
		return nil, err
	}

	keywordResults, err := r.KeywordSearch(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	return rrfFusion(semanticResults, keywordResults, opts.Limit), nil
}

// rrfFusion combines two ranked lists using Reciprocal Rank Fusion.
// score = 1/(k + rank_semantic) + 1/(k + rank_keyword)
func rrfFusion(semantic, keyword []memory.ScoredMemory, limit int) []memory.ScoredMemory {
	scores := make(map[uuid.UUID]float32)
	memMap := make(map[uuid.UUID]memory.ScoredMemory)

	for rank, sm := range semantic {
		scores[sm.Memory.ID] += 1.0 / float32(rrfK+rank+1)
		memMap[sm.Memory.ID] = sm
	}

	for rank, sm := range keyword {
		scores[sm.Memory.ID] += 1.0 / float32(rrfK+rank+1)
		if _, exists := memMap[sm.Memory.ID]; !exists {
			memMap[sm.Memory.ID] = sm
		}
	}

	results := make([]memory.ScoredMemory, 0, len(memMap))
	for id, sm := range memMap {
		results = append(results, memory.ScoredMemory{
			Memory:      sm.Memory,
			Score:       scores[id],
			MatchSource: "hybrid",
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	if limit <= 0 {
		limit = defaultLimit
	}
	if limit > len(results) {
		limit = len(results)
	}
	return results[:limit]
}

// filter params builders

func buildSemanticParams(embedding []float32, opts search.SearchFilters) sqlcgen.SemanticSearchParams {
	limit := int32(opts.Limit)
	if limit <= 0 {
		limit = int32(defaultLimit)
	}

	return sqlcgen.SemanticSearchParams{
		Embedding:          pgvector.NewVector(embedding),
		TenantID:           opts.TenantID,
		FilterUser:         opts.UserID != uuid.Nil,
		UserID:             opts.UserID,
		FilterTypes:        len(opts.MemoryTypes) > 0,
		MemoryTypes:        memoryTypesToStrings(opts.MemoryTypes),
		FilterSources:      len(opts.Sources) > 0,
		Sources:            opts.Sources,
		FilterContentTypes: len(opts.ContentTypes) > 0,
		ContentTypes:       contentTypesToStrings(opts.ContentTypes),
		FilterStartDate:    !opts.StartDate.IsZero(),
		StartDate:          toTimestampz(opts.StartDate),
		FilterEndDate:      !opts.EndDate.IsZero(),
		EndDate:            toTimestampz(opts.EndDate),
		ResultLimit:        limit,
	}
}

func buildKeywordParams(query string, opts search.SearchFilters) sqlcgen.KeywordSearchParams {
	limit := int32(opts.Limit)
	if limit <= 0 {
		limit = int32(defaultLimit)
	}

	return sqlcgen.KeywordSearchParams{
		Query:              query,
		TenantID:           opts.TenantID,
		FilterUser:         opts.UserID != uuid.Nil,
		UserID:             opts.UserID,
		FilterTypes:        len(opts.MemoryTypes) > 0,
		MemoryTypes:        memoryTypesToStrings(opts.MemoryTypes),
		FilterSources:      len(opts.Sources) > 0,
		Sources:            opts.Sources,
		FilterContentTypes: len(opts.ContentTypes) > 0,
		ContentTypes:       contentTypesToStrings(opts.ContentTypes),
		FilterStartDate:    !opts.StartDate.IsZero(),
		StartDate:          toTimestampz(opts.StartDate),
		FilterEndDate:      !opts.EndDate.IsZero(),
		EndDate:            toTimestampz(opts.EndDate),
		ResultLimit:        limit,
	}
}

// type converters

func memoryTypesToStrings(types []memory.MemoryType) []string {
	s := make([]string, len(types))
	for i, t := range types {
		s[i] = string(t)
	}
	return s
}

func contentTypesToStrings(types []memory.ContentType) []string {
	s := make([]string, len(types))
	for i, t := range types {
		s[i] = string(t)
	}
	return s
}

// row → domain mappers

func semanticRowToDomain(row sqlcgen.SemanticSearchRow) (*memory.Memory, error) {
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

func keywordRowToDomain(row sqlcgen.KeywordSearchRow) (*memory.Memory, error) {
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
