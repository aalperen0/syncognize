package memory

import "errors"

var (
	ErrNotFound      = errors.New("memory not found")
	ErrDuplicateHash = errors.New("duplicate content hash")
)
