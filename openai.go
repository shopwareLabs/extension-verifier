package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type newCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type openaiResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int    `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Message struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"message"`
		Index        int    `json:"index"`
		Logprobs     any    `json:"logprobs"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}

func NewCompletionRequest(text string) (string, error) {
	req := newCompletionRequest{
		Model: "phi-4",
		Messages: []Message{
			{
				Role:    "user",
				Content: text,
			},
		},
	}
	reqBody, err := json.Marshal(req)

	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", "http://localhost:1234/v1/chat/completions", bytes.NewBuffer(reqBody))

	if err != nil {
		return "", err
	}

	r.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(r)

	if err != nil {
		return "", err
	}

	content, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	var openaiResp openaiResponse

	if err := json.Unmarshal(content, &openaiResp); err != nil {
		return "", err
	}

	if len(openaiResp.Choices) == 0 {
		fmt.Println(string(content))
		return "", fmt.Errorf("no choices in response")
	}

	return openaiResp.Choices[0].Message.Content, nil
}
