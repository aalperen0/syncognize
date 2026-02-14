package search

import (
	"time"

	"github.com/aalperen0/syncognize/internal/domain/memory"
	"github.com/google/uuid"
)

type SearchFilters struct {
	TenantID     uuid.UUID
	UserID       uuid.UUID
	MemoryTypes  []memory.MemoryType
	Sources      []string
	ContentTypes []memory.ContentType
	StartDate    time.Time
	EndDate      time.Time
	Limit        int
}
