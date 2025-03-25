package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/sync/errgroup"
)

type StylintError struct {
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Rule     string `json:"rule"`
	Severity string `json:"severity"`
	Text     string `json:"text"`
}

type StylelintOutput []struct {
	Source                string         `json:"source"`
	Deprecations          []StylintError `json:"deprecations"`
	InvalidOptionWarnings []any          `json:"invalidOptionWarnings"`
	ParseErrors           []any          `json:"parseErrors"`
	Errored               bool           `json:"errored"`
	Warnings              []StylintError `json:"warnings"`
}

type StyleLint struct{}

func (s StyleLint) Name() string {
	return "stylelint"
}

func (s StyleLint) Check(ctx context.Context, check *Check, config ToolConfig) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	paths := append(config.StorefrontDirectories, config.AdminDirectories...)
	var gr errgroup.Group

	for _, p := range paths {
		p := p

		if !path.IsAbs(p) {
			p = path.Join(cwd, p)
		}

		hasSCSS, err := hasSCSSFiles(p)
		if err != nil {
			return err
		}

		if !hasSCSS {
			continue
		}

		gr.Go(func() error {
			stylelint := exec.CommandContext(ctx, "node", path.Join(cwd, "tools", "js", "node_modules", ".bin", "stylelint"),
				"--formatter=json",
				"--config", path.Join(cwd, "tools", "js", fmt.Sprintf("stylelint.config.%s.mjs", path.Base(p))),
				"--ignore-pattern", "dist/**",
				"--ignore-pattern", "vendor/**",
				fmt.Sprintf("%s/**/*.scss", p),
			)
			stylelint.Dir = p

			log, _ := stylelint.CombinedOutput()

			var stylelintOutput StylelintOutput

			if err := json.Unmarshal(log, &stylelintOutput); err != nil {
				return fmt.Errorf("failed to unmarshal stylelint output: %w, %s", err, string(log))
			}

			for _, diagnostic := range stylelintOutput {
				fixedPath := strings.TrimPrefix(strings.TrimPrefix(diagnostic.Source, "/private"), config.RootDir+"/")

				for _, msg := range diagnostic.Warnings {
					check.AddResult(CheckResult{
						Path:       fixedPath,
						Line:       msg.Line,
						Message:    msg.Text,
						Severity:   msg.Severity,
						Identifier: fmt.Sprintf("stylelint/%s", msg.Rule),
					})
				}

				for _, msg := range diagnostic.Deprecations {
					check.AddResult(CheckResult{
						Path:       fixedPath,
						Line:       msg.Line,
						Message:    msg.Text,
						Severity:   msg.Severity,
						Identifier: fmt.Sprintf("stylelint/%s", msg.Rule),
					})
				}
			}

			return nil
		})
	}

	return gr.Wait()
}

func (s StyleLint) Fix(ctx context.Context, config ToolConfig) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	paths := append(config.StorefrontDirectories, config.AdminDirectories...)

	var gr errgroup.Group

	for _, p := range paths {
		p := p

		if !path.IsAbs(p) {
			p = path.Join(cwd, p)
		}

		hasSCSS, err := hasSCSSFiles(p)
		if err != nil {
			return err
		}

		if !hasSCSS {
			continue
		}

		gr.Go(func() error {
			stylelint := exec.CommandContext(ctx, "node", path.Join(cwd, "tools", "js", "node_modules", ".bin", "stylelint"),
				"--config", path.Join(cwd, "tools", "js", fmt.Sprintf("stylelint.config.%s.mjs", path.Base(p))),
				"--ignore-pattern", "dist/**",
				"--ignore-pattern", "vendor/**",
				"**/*.scss",
				"--fix",
			)
			stylelint.Dir = p
			stylelint.Stdout = os.Stdout
			stylelint.Stderr = os.Stderr

			return stylelint.Run()
		})
	}

	return gr.Wait()
}

func (s StyleLint) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}

func init() {
	AddTool(StyleLint{})
}

func hasSCSSFiles(dir string) (bool, error) {
	matches, err := filepath.Glob(filepath.Join(dir, "**", "*.scss"))
	if err != nil {
		return false, err
	}

	return len(matches) > 0, nil
}
