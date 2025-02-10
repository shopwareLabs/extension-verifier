package main

import (
	"context"
	"fmt"

	"github.com/shopware/shopware-cli/extension"
)

type SWCLI struct{}

func (s SWCLI) Check(ctx context.Context, check *Check, config ToolConfig) error {
	validationContext := extension.ValidationContext{Extension: config.Extension}

	config.Extension.Validate(ctx, &validationContext)

	for _, err := range validationContext.Errors() {
		check.AddResult(CheckResult{
			Path:       "",
			Line:       0,
			Message:    err.Message,
			Identifier: fmt.Sprintf("shopware-cli/%s", err.Identifier),
			Severity:   "error",
		})
	}

	for _, err := range validationContext.Warnings() {
		check.AddResult(CheckResult{
			Path:       "",
			Line:       0,
			Message:    err.Message,
			Identifier: fmt.Sprintf("shopware-cli/%s", err.Identifier),
			Severity:   "warning",
		})
	}

	return nil
}

func (s SWCLI) Fix(ctx context.Context, config ToolConfig) error {
	return nil
}

func (s SWCLI) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}
