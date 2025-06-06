package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/shopware/extension-verifier/internal/llm"
	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/extension-verifier/internal/twig"
	"github.com/shopware/shopware-cli/extension"
	"github.com/shopware/shopware-cli/logging"
	"github.com/spf13/cobra"
)

const systemPrompt = `
You are a helper agent to help to upgrade Twig templates. I will give you the old and new template happend in the Software and as third the extended template. Apply the changes happen between old and new template to the extended template.
- Do only the necessary changes to the extended template.
- Do only modify the content inside the block and dont add new blocks
- Please also only output the modified extended template nothing more.
- Adjust also HTML elements to be more accessibility friendly.
- If in a {% block %} is {{ parent() }}, ignore it and dont modify the content of the block
`

var twigUpgradeCommand = &cobra.Command{
	Use:   "twig-upgrade [path] [old-shopware-version] [new-shopware-version]",
	Short: "Experimental upgrade of Twig templates using AI",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Error("Extension Verifier project have been moved into shopware-cli itself")
		log.Error("The new command is: docker run --rm -v $(pwd):/ext shopware/shopware-cli shopware-cli extension ai twig-upgrade")
		log.Error("Sleeping for 30 seconds before running the old command")
		time.Sleep(30 * time.Second)

		ext, err := extension.GetExtensionByFolder(args[0])

		if err != nil {
			return err
		}

		logging.FromContext(cmd.Context()).Debugf("System Prompt:\n%s\n", systemPrompt)

		toolCfg, err := tool.ConvertExtensionToToolConfig(ext)

		if err != nil {
			return err
		}

		client, err := llm.NewLLMClient(cmd.Flag("provider").Value.String())

		if err != nil {
			return err
		}

		options := &llm.LLMOptions{
			Model:        cmd.Flag("model").Value.String(),
			SystemPrompt: systemPrompt,
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

			defer func() {
				if err := os.RemoveAll(oldVersion); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove old version directory: %v\n", err)
				}
			}()

			newVersion, err := cloneShopwareStorefront(args[2])

			if err != nil {
				return err
			}

			defer func() {
				if err := os.RemoveAll(newVersion); err != nil {
					fmt.Fprintf(os.Stderr, "Failed to remove new version directory: %v\n", err)
				}
			}()

			err = filepath.Walk(twigFolder, func(file string, info os.FileInfo, _ error) error {
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

				log.Info("Processing file", "file", file)

				logging.FromContext(cmd.Context()).Debugf("Input to LLM for file %s:\n%s\n", file, str.String())

				text, err := client.Generate(cmd.Context(), str.String(), options)

				if err != nil {
					return err
				}

				if strings.Contains(text, "</think>") {
					thinkEndIndex := strings.Index(text, "</think>")
					if thinkEndIndex != -1 {
						text = text[thinkEndIndex+len("</think>"):]
					}
				}

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

			if err != nil {
				return err
			}
		}

		return nil
	},
}

func cloneShopwareStorefront(version string) (string, error) {
	tempDir, err := os.MkdirTemp(os.TempDir(), "shopware")

	if err != nil {
		return "", err
	}

	git := exec.Command("git", "-c", "advice.detachedHead=false", "clone", "-q", "--branch", "v"+version, "https://github.com/shopware/storefront", tempDir, "--depth", "1")
	output, err := git.CombinedOutput()

	if err != nil {
		log.Print(string(output))
	}

	return tempDir, err
}

func init() {
	twigUpgradeCommand.Flags().String("model", "gemma3:4b", "The model to use for the upgrade")
	twigUpgradeCommand.Flags().String("provider", "ollama", "The provider to use for the upgrade")
	rootCmd.AddCommand(twigUpgradeCommand)
}
