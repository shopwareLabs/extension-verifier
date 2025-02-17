package tool

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/shopware/extension-verifier/internal/admintwiglinter"
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type AdminTwigLinter struct{}

func (a AdminTwigLinter) Check(ctx context.Context, check *Check, config ToolConfig) error {
	fixers := admintwiglinter.GetFixers(version.Must(version.NewVersion(config.MinShopwareVersion)))

	for _, p := range GetAdminFolders(config) {
		err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".twig" {
				return nil
			}

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			parsed, err := html.NewParser(string(file))

			if err != nil {
				return fmt.Errorf("failed to parse %s: %w", path, err)
			}

			for _, fixer := range fixers {
				for _, message := range fixer.Check(parsed) {
					check.AddResult(CheckResult{
						Message:    message.Message,
						Path:       strings.TrimPrefix(strings.TrimPrefix(path, "/private"), config.Extension.GetPath()+"/"),
						Line:       0,
						Severity:   message.Severity,
						Identifier: fmt.Sprintf("admintwiglinter/%s", message.Identifier),
					})
				}
			}

			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (a AdminTwigLinter) Fix(ctx context.Context, config ToolConfig) error {
	fixers := admintwiglinter.GetFixers(version.Must(version.NewVersion(config.MinShopwareVersion)))

	for _, p := range GetAdminFolders(config) {
		err := filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".twig" {
				return nil
			}

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			parsed, err := html.NewParser(string(file))

			if err != nil {
				return err
			}

			for _, fixer := range fixers {
				if err := fixer.Fix(parsed); err != nil {
					return err
				}
			}

			var buf strings.Builder

			for _, node := range parsed {
				buf.WriteString(node.Dump())
			}

			return os.WriteFile(path, []byte(buf.String()), os.ModePerm)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (a AdminTwigLinter) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}

func init() {
	AddTool(AdminTwigLinter{})
}
