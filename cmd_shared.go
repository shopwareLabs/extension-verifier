package main

import (
	"fmt"
	"strings"

	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/shopware-cli/extension"
)

func filterTools(tools []tool.Tool, only string) ([]tool.Tool, error) {
	if only == "" {
		return tools, nil
	}

	var filteredTools []tool.Tool
	requestedTools := strings.Split(only, ",")

	for _, requestedTool := range requestedTools {
		requestedTool = strings.TrimSpace(requestedTool)
		found := false

		for _, t := range tools {
			if t.Name() == requestedTool {
				filteredTools = append(filteredTools, t)
				found = true
				break
			}
		}

		if !found {
			return nil, fmt.Errorf("tool with name %q not found", requestedTool)
		}
	}

	return filteredTools, nil
}

func getToolConfig(path string) (*tool.ToolConfig, error) {
	var toolCfg *tool.ToolConfig
	var err error

	if tool.IsProject(path) {
		toolCfg, err = tool.GetConfigFromProject(path)
		if err != nil {
			return nil, err
		}
	} else {
		ext, err := extension.GetExtensionByFolder(path)
		if err != nil {
			return nil, err
		}

		toolCfg, err = tool.ConvertExtensionToToolConfig(ext)
		if err != nil {
			return nil, err
		}
	}

	return toolCfg, nil
}
