package main

import (
	"fmt"
	"strings"

	"github.com/shopware/extension-verifier/internal/tool"
)

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
