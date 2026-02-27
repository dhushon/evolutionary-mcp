package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPMLClient is an HTTP implementation of the MLClient interface.
type HTTPMLClient struct {
	url string
}

// NewHTTPMLClient creates a new HTTPMLClient.
func NewHTTPMLClient(url string) *HTTPMLClient {
	return &HTTPMLClient{url: url}
}

// GetEmbedding returns the embedding for a given text.
func (c *HTTPMLClient) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	requestBody, err := json.Marshal(map[string]string{"text": text})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.url+"/embedding", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get embedding: status code %d", resp.StatusCode)
	}

	var embedding []float32
	if err := json.NewDecoder(resp.Body).Decode(&embedding); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return embedding, nil
}
