package main

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

//go:embed configs/phpstan.neon.sw6
var phpstanConfigSW6 []byte

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
	Errors []any `json:"errors"`
}

type PhpStan struct{}

func (p PhpStan) Check(ctx context.Context, check *Check, config ToolConfig) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if err := installComposerDeps(config.RootDir, config.CheckAgainst); err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(config.RootDir, "phpstan.neon"), phpstanConfigSW6, 0644); err != nil {
		return err
	}

	phpstan := exec.CommandContext(ctx, "php", "-dmemory_limit=2G", path.Join(cwd, "tools", "phpstan", "vendor", "bin", "phpstan"), "analyse", "--no-progress", "--no-interaction", "--error-format=json")
	phpstan.Dir = config.RootDir

	log, _ := phpstan.Output()

	var phpstanResult PhpStanOutput

	if err := json.Unmarshal(log, &phpstanResult); err != nil {
		if strings.Contains(err.Error(), "cannot unmarshal array into Go struct field PhpStanOutput.file") {
			return nil
		}

		return fmt.Errorf("failed to unmarshal phpstan output: %w", err)
	}

	for fileName, file := range phpstanResult.Files {
		for _, message := range file.Messages {
			check.AddResult(CheckResult{
				Path:       strings.TrimPrefix(strings.TrimPrefix(fileName, "/private"), config.RootDir+"/"),
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
