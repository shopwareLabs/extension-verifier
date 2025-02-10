package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"golang.org/x/sync/errgroup"
)

type EslintOutput []struct {
	FilePath string `json:"filePath"`
	Messages []struct {
		RuleID    string `json:"ruleId"`
		Severity  int    `json:"severity"`
		Message   string `json:"message"`
		Line      int    `json:"line"`
		Column    int    `json:"column"`
		NodeType  string `json:"nodeType"`
		EndLine   int    `json:"endLine"`
		EndColumn int    `json:"endColumn"`
		Fix       struct {
			Range []int  `json:"range"`
			Text  string `json:"text"`
		} `json:"fix,omitempty"`
		MessageID string `json:"messageId,omitempty"`
	} `json:"messages"`
	SuppressedMessages  []any  `json:"suppressedMessages"`
	ErrorCount          int    `json:"errorCount"`
	FatalErrorCount     int    `json:"fatalErrorCount"`
	WarningCount        int    `json:"warningCount"`
	FixableErrorCount   int    `json:"fixableErrorCount"`
	FixableWarningCount int    `json:"fixableWarningCount"`
	Source              string `json:"source"`
	UsedDeprecatedRules []any  `json:"usedDeprecatedRules"`
}

type Eslint struct{}

func (e Eslint) Check(ctx context.Context, check *Check, config ToolConfig) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	paths := getStorefrontPaths(config)

	var gr errgroup.Group

	for _, p := range paths {
		p := p
		gr.Go(func() error {
			eslint := exec.CommandContext(ctx, "node", path.Join(cwd, "tools", "eslint", "node_modules", ".bin", "eslint"), "--format=json", "--config", path.Join(cwd, "tools", "eslint", path.Base(p)+".config.mjs"), "--ignore-pattern", "dist/**", "--ignore-pattern", "vendor/**")
			eslint.Dir = p

			log, _ := eslint.CombinedOutput()

			var eslintOutput EslintOutput

			if err := json.Unmarshal(log, &eslintOutput); err != nil {
				return fmt.Errorf("failed to unmarshal eslint output: %w, %s", err, string(log))
			}

			for _, diagnostic := range eslintOutput {
				fixedPath := strings.TrimPrefix(strings.TrimPrefix(diagnostic.FilePath, "/private"), config.Extension.GetPath()+"/")

				for _, message := range diagnostic.Messages {
					severity := "warn"

					if message.Severity == 2 {
						severity = "error"
					}

					check.AddResult(CheckResult{
						Path:       fixedPath,
						Line:       message.Line,
						Message:    message.Message,
						Severity:   severity,
						Identifier: fmt.Sprintf("eslint/%s", message.RuleID),
					})
				}
			}

			return nil
		})
	}

	return gr.Wait()
}

func (e Eslint) Fix(ctx context.Context, config ToolConfig) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	paths := getStorefrontPaths(config)

	var gr errgroup.Group

	for _, p := range paths {
		p := p
		gr.Go(func() error {
			eslint := exec.CommandContext(ctx, "node", path.Join(cwd, "tools", "eslint", "node_modules", ".bin", "eslint"), "--config", path.Join(cwd, "tools", "eslint", path.Base(p)+".config.mjs"), "--ignore-pattern", "dist/**", "--ignore-pattern", "vendor/**", "--fix")
			eslint.Dir = p

			log, _ := eslint.CombinedOutput()

			fmt.Println(string(log))

			return nil
		})
	}

	return gr.Wait()
}

func (e Eslint) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}
