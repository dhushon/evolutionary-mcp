package services

import "context"

// MLClient is an interface for communicating with the ML sidecar.
type MLClient interface {
	// GetEmbedding returns the embedding for a given text.
	GetEmbedding(ctx context.Context, text string) ([]float32, error)
}
