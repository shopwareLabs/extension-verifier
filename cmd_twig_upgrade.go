package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/extension-verifier/internal/twig"
	"github.com/shopware/shopware-cli/extension"
	"github.com/spf13/cobra"
	"google.golang.org/api/option"
)

var twigUpgradeCommand = &cobra.Command{
	Use:   "twig-upgrade [path] [old-shopware-version] [new-shopware-version]",
	Short: "Experimental upgrade of Twig templates using AI",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ext, err := extension.GetExtensionByFolder(args[0])

		if err != nil {
			return err
		}

		toolCfg, err := tool.ConvertExtensionToToolConfig(ext)

		if err != nil {
			return err
		}

		apiKey := os.Getenv("GEMINI_API_KEY")

		if apiKey == "" {
			return fmt.Errorf("GEMINI_API_KEY is not set")
		}

		client, err := genai.NewClient(cmd.Context(), option.WithAPIKey(apiKey))

		if err != nil {
			return err
		}

		for _, sourceDirectory := range toolCfg.SourceDirectories {
			twigFolder := path.Join(sourceDirectory, "Resources", "views", "storefront")

			if _, err := os.Stat(twigFolder); os.IsNotExist(err) {
				return nil
			}

			oldVersion, err := cloneShopwareStorefront(args[1])

			if err != nil {
				return err
			}

			newVersion, err := cloneShopwareStorefront(args[2])

			if err != nil {
				return err
			}

			defer func() {
				if err := os.RemoveAll(oldVersion); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove old version directory: %v\n", err)
				}
			}()
			defer func() {
				if err := os.RemoveAll(newVersion); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove new version directory: %v\n", err)
				}
			}()

			_ = filepath.Walk(twigFolder, func(file string, info os.FileInfo, _ error) error {
				if info.IsDir() {
					return nil
				}

				if filepath.Ext(file) != ".twig" {
					return nil
				}

				content, err := os.ReadFile(file)

				if err != nil {
					return err
				}

				ast, err := twig.ParseTemplate(string(content))

				if err != nil {
					return err
				}

				extends := ast.Extends()

				if extends == nil {
					return nil
				}

				tpl := extends.Template

				if tpl[0] == '@' {
					tplParts := strings.Split(tpl, "/")
					tplParts = tplParts[1:]
					tpl = strings.Join(tplParts, "/")
				}

				oldTemplateText, err := os.ReadFile(path.Join(oldVersion, "Resources", "views", tpl))

				if err != nil {
					fmt.Printf("Template %s not found in old version\n", tpl)
					return nil
				}

				newTemplateText, err := os.ReadFile(path.Join(newVersion, "Resources", "views", tpl))

				if err != nil {
					fmt.Printf("Template %s not found in new version\n", tpl)
					return nil
				}

				var str strings.Builder
				str.WriteString("You are a helper agent to help to upgrade Twig templates. I will give you the old and new template happend in the Software and as third the extended template. Apply the changes happen between old and new template to the extended template.\n")
				str.WriteString("Follow following rules while making adjustments to the extended template:\n")
				str.WriteString("- Do only the necessary changes to the extended template.\n")
				str.WriteString("- Do only modify the content inside the block and dont add new blocks\n")
				str.WriteString("- Please also only output the modified extended template nothing more.\n")
				str.WriteString("- Adjust also HTML elements to be more accessibility friendly.\n")
				str.WriteString("- If in a {% block %} is {{ parent() }}, ignore it and dont modify the content of the block\n")
				str.WriteString("\n")
				str.WriteString("This was the old template:\n")
				str.WriteString("```twig\n")
				str.WriteString(string(oldTemplateText))
				str.WriteString("\n```\n")
				str.WriteString("and this is the new one:\n")
				str.WriteString("```twig\n")
				str.WriteString(string(newTemplateText))
				str.WriteString("\n```\n")
				str.WriteString("and this is my template:\n")
				str.WriteString("```twig\n")
				str.WriteString(string(content))
				str.WriteString("\n```")

				resp, err := generateContent(cmd.Context(), client, str.String())

				if err != nil {
					return err
				}

				text := string(resp.Candidates[0].Content.Parts[0].(genai.Text))

				start := strings.Index(text, "```twig")
				end := strings.LastIndex(text, "```")

				if start == -1 || end == -1 {
					return nil
				}

				text = strings.TrimPrefix(text[start+7:end], "\n")

				contentStr := string(content)
				if strings.TrimSpace(text) == strings.TrimSpace(contentStr) {
					return nil
				}

				return os.WriteFile(file, []byte(text), os.ModePerm)
			})
		}
		return nil
	},
}

func generateContent(ctx context.Context, client *genai.Client, message string) (*genai.GenerateContentResponse, error) {
	resp, err := client.GenerativeModel("gemini-2.0-pro-exp-02-05").GenerateContent(ctx, genai.Text(message))

	if err != nil {
		if strings.Contains(err.Error(), "Resource has been exhausted") {
			fmt.Println("Resource exhausted, waiting 15 seconds before retrying")
			time.Sleep(15 * time.Second)

			return generateContent(ctx, client, message)
		}
	}

	return resp, err
}

func cloneShopwareStorefront(version string) (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "shopware")

	if err != nil {
		return "", err
	}

	git := exec.Command("git", "clone", "--branch", "v"+version, "https://github.com/shopware/storefront", tempDir, "--depth", "1")
	git.Stdout = os.Stdout
	git.Stderr = os.Stderr

	if err := git.Run(); err != nil {
		return "", err
	}

	return tempDir, nil
}

func init() {
	rootCmd.AddCommand(twigUpgradeCommand)
}
