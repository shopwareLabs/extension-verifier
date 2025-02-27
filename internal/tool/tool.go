package tool

import (
	"context"
	"os"
	"path"

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
	MinShopwareVersion string
	MaxShopwareVersion string
	CheckAgainst       string
	Extension          extension.Extension
}

type Tool interface {
	Check(ctx context.Context, check *Check, config ToolConfig) error
	Fix(ctx context.Context, config ToolConfig) error
	Format(ctx context.Context, config ToolConfig, dryRun bool) error
}

func GetAdminFolders(config ToolConfig) []string {
	paths := []string{
		path.Join(config.Extension.GetResourcesDir(), "app", "administration"),
	}

	for _, bundle := range config.Extension.GetExtensionConfig().Build.ExtraBundles {
		paths = append(paths, path.Join(config.Extension.GetRootDir(), bundle.Path, "Resources", "app", "administration"))
	}

	filteredPaths := make([]string, 0)
	for _, p := range paths {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			filteredPaths = append(filteredPaths, p)
		}
	}

	paths = filteredPaths

	return paths
}

func GetJSFolders(config ToolConfig) []string {
	paths := []string{
		path.Join(config.Extension.GetResourcesDir(), "app", "storefront"),
		path.Join(config.Extension.GetResourcesDir(), "app", "administration"),
	}

	for _, bundle := range config.Extension.GetExtensionConfig().Build.ExtraBundles {
		paths = append(paths, path.Join(config.Extension.GetRootDir(), bundle.Path, "Resources", "app", "storefront"))
		paths = append(paths, path.Join(config.Extension.GetRootDir(), bundle.Path, "Resources", "app", "administration"))
	}

	filteredPaths := make([]string, 0)
	for _, p := range paths {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			filteredPaths = append(filteredPaths, p)
		}
	}

	paths = filteredPaths

	return paths
}
