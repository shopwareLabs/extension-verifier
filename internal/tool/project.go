package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"slices"
	"strings"

	"github.com/shopware/shopware-cli/extension"
	"github.com/shopware/shopware-cli/shop"
	"github.com/shyim/go-version"
)

func IsProject(root string) bool {
	composerJson := path.Join(root, "composer.json")

	if _, err := os.Stat(composerJson); os.IsNotExist(err) {
		return false
	}

	var composerJsonData struct {
		Type string `json:"type"`
	}

	file, err := os.Open(composerJson)

	if err != nil {
		return false
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close composer.json: %w", closeErr)
		}
	}()

	if err := json.NewDecoder(file).Decode(&composerJsonData); err != nil {
		return false
	}

	return composerJsonData.Type == "project"
}

func getShopwareConstraint(root string) (*version.Constraints, error) {
	composerJson := path.Join(root, "composer.json")

	file, err := os.Open(composerJson)

	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close composer.json: %w", closeErr)
		}
	}()

	var composerJsonData struct {
		Require struct {
			Shopware string `json:"shopware/core"`
		} `json:"require"`
	}

	if err := json.NewDecoder(file).Decode(&composerJsonData); err != nil {
		return nil, err
	}

	if composerJsonData.Require.Shopware == "" {
		return nil, fmt.Errorf("shopware/core is not required")
	}

	cst, err := version.NewConstraint(composerJsonData.Require.Shopware)

	if err != nil {
		return nil, err
	}

	return &cst, nil
}

func GetConfigFromProject(root string) (*ToolConfig, error) {
	constraint, err := getShopwareConstraint(root)

	if err != nil {
		return nil, err
	}

	extensions := extension.FindExtensionsFromProject(context.Background(), root)

	sourceDirectories := []string{}
	adminDirectories := []string{}
	storefrontDirectories := []string{}

	vendorPath := path.Join(root, "vendor")

	shopCfg, err := shop.ReadConfig(path.Join(root, ".shopware-project.yml"), true)

	if err != nil {
		return nil, err
	}

	excludeExtensions := []string{}

	if shopCfg.Validation != nil {
		for _, ignore := range shopCfg.Validation.IgnoreExtensions {
			excludeExtensions = append(excludeExtensions, ignore.Name)
		}
	}

	for _, ext := range extensions {
		extName, err := ext.GetName()

		if err != nil {
			return nil, err
		}

		// Skip plugins in vendor folder
		if strings.HasPrefix(ext.GetRootDir(), vendorPath) || slices.Contains(excludeExtensions, extName) {
			continue
		}

		sourceDirectories = append(sourceDirectories, ext.GetRootDir())
		adminDirectories = append(adminDirectories, getAdminFolders(ext)...)
		storefrontDirectories = append(storefrontDirectories, getStorefrontFolders(ext)...)
	}

	var validationIgnores []ToolConfigIgnore

	if shopCfg.Validation != nil {
		for _, ignore := range shopCfg.Validation.Ignore {
			validationIgnores = append(validationIgnores, ToolConfigIgnore{
				Identifier: ignore.Identifier,
				Path:       ignore.Path,
				Message:    ignore.Message,
			})
		}
	}

	toolCfg := &ToolConfig{
		RootDir:               root,
		SourceDirectories:     sourceDirectories,
		AdminDirectories:      adminDirectories,
		StorefrontDirectories: storefrontDirectories,
		ValidationIgnores:     validationIgnores,
	}

	if err := determineVersionRange(toolCfg, constraint); err != nil {
		return nil, err
	}

	return toolCfg, nil
}
