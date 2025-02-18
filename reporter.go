package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/shopware/extension-verifier/internal/tool"
)

func doCIReport(result *tool.Check) error {
	isGitHubAction := os.Getenv("GITHUB_ACTIONS") == "true"

	if isGitHubAction {
		stepSummary := os.Getenv("GITHUB_STEP_SUMMARY")

		if stepSummary != "" {
			if err := os.WriteFile(stepSummary, []byte(convertResultsToMarkdown(result.Results)), 0644); err != nil {
				return fmt.Errorf("failed to write step summary: %w", err)
			}
		}

		for _, res := range result.Results {
			if res.Line == 0 {
				fmt.Printf("::%s file=%s::%s\n", res.Severity, res.Path, res.Message)
			} else {
				fmt.Printf("::%s file=%s,line=%d::%s\n", res.Severity, res.Path, res.Line, res.Message)
			}
		}
	}

	return nil

}

func convertResultsToMarkdown(check []tool.CheckResult) string {
	var builder strings.Builder

	builder.WriteString("# Results\n\n")

	builder.WriteString("| Severity | Identifier | File | Message | \n")
	builder.WriteString("| --- | --- | --- | --- |\n")

	for _, result := range check {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", result.Severity, result.Identifier, result.Path, result.Message))
	}

	builder.WriteString("\n")

	return builder.String()
}
