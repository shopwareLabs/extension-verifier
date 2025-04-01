package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// Client represents an Ollama API client
type Client struct {
	host   string
	model  string
	client *http.Client
}

// GenerateRequest represents the request body for text generation
type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	System string `json:"system,omitempty"`
}

// GenerateResponse represents the response from the generate endpoint
type GenerateResponse struct {
	Model           string `json:"model"`
	CreatedAt       string `json:"created_at"`
	Response        string `json:"response"`
	TotalDuration   int64  `json:"total_duration"`
	LoadDuration    int64  `json:"load_duration"`
	PromptEvalCount int    `json:"prompt_eval_count"`
	EvalCount       int    `json:"eval_count"`
	EvalDuration    int64  `json:"eval_duration"`
}

// newOllamaClient creates a new Ollama client instance
func newOllamaClient() *Client {
	host := os.Getenv("OLLAMA_HOST")
	if host == "" {
		host = "http://localhost:11434"
	}

	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "deepseek-r1:1.5b"
	}

	return &Client{
		host:  host,
		model: model,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Generate sends a generation request to the Ollama API
func (c *Client) Generate(ctx context.Context, prompt string, options *LLMOptions) (string, error) {
	reqBody := GenerateRequest{
		Model:  options.Model,
		Prompt: prompt,
		Stream: false,
	}

	if options.SystemPrompt != "" {
		reqBody.System = options.SystemPrompt
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/generate", c.host), bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var response GenerateResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Response, nil
}

// SetModel changes the model used by the client
func (c *Client) SetModel(model string) {
	c.model = model
}

// SetHost changes the host used by the client
func (c *Client) SetHost(host string) {
	c.host = host
}
