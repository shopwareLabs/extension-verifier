package llm

import "context"

type LLMOptions struct {
	Model        string
	SystemPrompt string
}

type LLMClient interface {
	Generate(ctx context.Context, prompt string, options *LLMOptions) (string, error)
}

func NewLLMClient(provider string) LLMClient {
	switch provider {
	case "ollama":
		return newOllamaClient()
	case "gemini":
		return newGeminiClient()
	}

	panic("Invalid provider: " + provider)
}
