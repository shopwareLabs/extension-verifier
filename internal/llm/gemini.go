package llm

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/log"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type GeminiClient struct {
	client *genai.Client
}

func newGeminiClient() *GeminiClient {
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))

	if err != nil {
		log.Fatal(err)
	}

	return &GeminiClient{client: client}
}

func (c *GeminiClient) Generate(ctx context.Context, prompt string, options *LLMOptions) (string, error) {
	resp, err := c.client.GenerativeModel("gemini-2.0-pro-exp-02-05").GenerateContent(ctx, genai.Text(prompt))

	if err != nil {
		if strings.Contains(err.Error(), "Resource has been exhausted") {
			log.Warn("Resource exhausted, waiting 15 seconds before retrying")
			time.Sleep(15 * time.Second)

			return c.Generate(ctx, prompt, options)
		}
	}

	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil
}
