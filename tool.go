package main

import (
	"context"

	"github.com/shopware/shopware-cli/extension"
)

var availableTools = []Tool{Eslint{}, PhpStan{}, Rector{}, SWCLI{}}

type ToolConfig struct {
	MinShopwareVersion string
	MaxShopwareVersion string
	CheckAgainst       string
	Extension          extension.Extension
}

type Tool interface {
	Check(ctx context.Context, check *Check, config ToolConfig) error
	Fix(ctx context.Context, config ToolConfig) error
}
