package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Embedder struct {
	Model string
}

func NewEmbedder() *Embedder {
	return &Embedder{Model: "nomic-embed-text"}
}

type EmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type EmbeddingResponse struct {
	Embedding []float64 `json:"embedding"` // Ollama returns float64
}

func (e *Embedder) GenerateEmbedding(text string) ([]float32, error) {
	reqBody := EmbeddingRequest{
		Model:  e.Model,
		Prompt: text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post("http://localhost:11434/api/embeddings", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to call ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama API returned status: %d", resp.StatusCode)
	}

	var embeddingResp EmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&embeddingResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert float64 to float32 for pgvector
	float32Embedding := make([]float32, len(embeddingResp.Embedding))
	for i, v := range embeddingResp.Embedding {
		float32Embedding[i] = float32(v)
	}

	return float32Embedding, nil
}
