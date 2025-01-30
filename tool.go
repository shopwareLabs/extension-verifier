package main

import "context"

var availableTools = []Tool{Eslint{}, PhpStan{}, Rector{}}

type ToolConfig struct {
	RootDir                   string
	ShopwareVersionConstraint string
	MinShopwareVersion        string
	MaxShopwareVersion        string
	CheckAgainst              string
}

type Tool interface {
	Check(ctx context.Context, check *Check, config ToolConfig) error
	Fix(ctx context.Context, config ToolConfig) error
}
