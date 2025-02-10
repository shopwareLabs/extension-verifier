package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/shopware/extension-verifier/internal/twig"
)

type Twig struct{}

func (t Twig) Check(ctx context.Context, check *Check, config ToolConfig) error {
	return nil
}

func (t Twig) Fix(ctx context.Context, config ToolConfig) error {
	twigFolder := path.Join(config.Extension.GetResourcesDir(), "views", "storefront")

	if _, err := os.Stat(twigFolder); os.IsNotExist(err) {
		return nil
	}

	oldVersion, err := cloneShopwareStorefront(config.MinShopwareVersion)

	if err != nil {
		return err
	}

	newVersion, err := cloneShopwareStorefront(config.MaxShopwareVersion)

	if err != nil {
		return err
	}

	defer os.RemoveAll(oldVersion)
	defer os.RemoveAll(newVersion)

	return filepath.Walk(twigFolder, func(file string, info os.FileInfo, _ error) error {
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
		str.WriteString("You are a helper agent to help to upgrade Twig templates. I will give you the old and new template happend in the Software and as third the extended template. Apply the changes happend between old and new template to the extended template. Do only the necessary changes to the extended template.")
		str.WriteString("If a block calls parent(), ignore that as a difference. Please also only output the modified extended template nothing more.")
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

		req := str.String()

		os.WriteFile("llm.txt", []byte(req), os.ModePerm)

		text, err := NewCompletionRequest(req)

		if err != nil {
			return err
		}

		start := strings.Index(text, "```twig")
		end := strings.LastIndex(text, "```")

		if start == -1 || end == -1 {
			return nil
		}

		text = strings.TrimPrefix(text[start+7:end], "\n")

		return os.WriteFile(file, []byte(text), os.ModePerm)
	})
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
