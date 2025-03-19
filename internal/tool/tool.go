package tool

import (
	"context"

	"github.com/shopware/shopware-cli/extension"
)

var availableTools = []Tool{}

func AddTool(tool Tool) {
	availableTools = append(availableTools, tool)
}

func GetTools() []Tool {
	return availableTools
}

type ToolConfig struct {
	// The minimum version of Shopware that is supported
	MinShopwareVersion string
	// The maximum version of Shopware that is supported
	MaxShopwareVersion string
	// The version of Shopware that is checked against
	CheckAgainst string
	// The root directory of the extension/project
	RootDir string
	// Contains a list of directories that are considered as source code
	SourceDirectories []string
	// Contains a list of identifiers that are ignored
	ValidationIgnores []ToolConfigIgnore
	// Contains a list of directories that are considered as admin code
	AdminDirectories []string
	// Contains a list of directories that are considered as storefront code
	StorefrontDirectories []string

	Extension extension.Extension
}

type ToolConfigIgnore struct {
	Identifier string
	Path       string
	Message    string
}

type Tool interface {
	Check(ctx context.Context, check *Check, config ToolConfig) error
	Fix(ctx context.Context, config ToolConfig) error
	Format(ctx context.Context, config ToolConfig, dryRun bool) error
}
