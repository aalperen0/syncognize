package ingestion

import "errors"

var (
	ErrJobNotFound = errors.New("ingestion job not found")
	ErrBufferEmpty = errors.New("buffer is empty")
	ErrJobFailed   = errors.New("ingestion job failed to complete")
)
