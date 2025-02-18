package tool

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

//go:embed phpstan.neon.sw6
var phpstanConfigSW6 []byte

var possiblePHPStanConfigs = []string{
	"phpstan.neon",
	"phpstan.neon.dist",
}

type PhpStanOutput struct {
	Totals struct {
		Errors     int `json:"errors"`
		FileErrors int `json:"file_errors"`
	} `json:"totals"`
	Files map[string]struct {
		Errors   int `json:"errors"`
		Messages []struct {
			Message    string `json:"message"`
			Line       int    `json:"line"`
			Ignorable  bool   `json:"ignorable"`
			Identifier string `json:"identifier"`
		} `json:"messages"`
	} `json:"files"`
	Errors []string `json:"errors"`
}

type PhpStan struct{}

func (p PhpStan) configExists(pluginPath string) bool {
	for _, config := range possiblePHPStanConfigs {
		if _, err := os.Stat(path.Join(pluginPath, config)); err == nil {
			return true
		}
	}

	return false
}

func (p PhpStan) Check(ctx context.Context, check *Check, config ToolConfig) error {
	if config.Extension.GetType() == "app" {
		return nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if err := installComposerDeps(config.Extension, config.CheckAgainst); err != nil {
		return err
	}

	if !p.configExists(config.Extension.GetPath()) {
		if err := os.WriteFile(path.Join(config.Extension.GetPath(), "phpstan.neon"), phpstanConfigSW6, 0644); err != nil {
			return err
		}
	}

	phpstan := exec.CommandContext(ctx, "php", "-dmemory_limit=2G", path.Join(cwd, "tools", "php", "vendor", "bin", "phpstan"), "analyse", "--no-progress", "--no-interaction", "--error-format=json")
	phpstan.Dir = config.Extension.GetPath()

	log, _ := phpstan.Output()

	log = []byte(strings.ReplaceAll(string(log), "\"files\":[]", "\"files\":{}"))

	var phpstanResult PhpStanOutput

	if err := json.Unmarshal(log, &phpstanResult); err != nil {
		return fmt.Errorf("failed to unmarshal phpstan output: %w", err)
	}

	for _, error := range phpstanResult.Errors {
		check.AddResult(CheckResult{
			Path:       "phpstan.neon",
			Message:    error,
			Severity:   "error",
			Line:       0,
			Identifier: "phpstan/error",
		})
	}

	for fileName, file := range phpstanResult.Files {
		for _, message := range file.Messages {
			check.AddResult(CheckResult{
				Path:       strings.TrimPrefix(strings.TrimPrefix(fileName, "/private"), config.Extension.GetPath()+"/"),
				Line:       message.Line,
				Message:    message.Message,
				Severity:   "error",
				Identifier: fmt.Sprintf("phpstan/%s", message.Identifier),
			})
		}
	}

	return nil
}

func (p PhpStan) Fix(ctx context.Context, config ToolConfig) error {
	return nil
}

func (p PhpStan) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}

func init() {
	AddTool(PhpStan{})
}
