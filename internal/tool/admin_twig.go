package tool

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"

	"github.com/shopware/extension-verifier/internal/admintwiglinter"
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type AdminTwigLinter struct{}

func (a AdminTwigLinter) Check(ctx context.Context, check *Check, config ToolConfig) error {
	fixers := admintwiglinter.GetFixers(version.Must(version.NewVersion(config.MinShopwareVersion)))

	for _, p := range config.AdminDirectories {
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
						Path:       strings.TrimPrefix(strings.TrimPrefix(path, "/private"), config.RootDir+"/"),
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

	for _, p := range config.AdminDirectories {
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

			return os.WriteFile(path, []byte(parsed.Dump(0)), os.ModePerm)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (a AdminTwigLinter) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	dmp := diffmatchpatch.New()

	for _, p := range config.AdminDirectories {
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

			if dryRun {
				diffs := dmp.DiffMain(string(file), parsed.Dump(0), false)

				fmt.Println(dmp.DiffPrettyText(diffs))

				return nil
			} else {
				return os.WriteFile(path, []byte(parsed.Dump(0)), os.ModePerm)
			}
		})

		if err != nil {
			return err
		}
	}

	return nil

}

func init() {
	AddTool(AdminTwigLinter{})
}
