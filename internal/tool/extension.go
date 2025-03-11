package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"sort"

	"github.com/shopware/shopware-cli/extension"
	"github.com/shopware/shopware-cli/version"
)

func ConvertExtensionToToolConfig(ext extension.Extension) (*ToolConfig, error) {
	cfg := &ToolConfig{
		Extension:             ext,
		ValidationIgnores:     ext.GetExtensionConfig().Validation.Ignore,
		RootDir:               ext.GetPath(),
		SourceDirectories:     []string{ext.GetRootDir()},
		AdminDirectories:      getAdminFolders(ext),
		StorefrontDirectories: getStorefrontFolders(ext),
	}

	if err := determineVersionRange(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func determineVersionRange(cfg *ToolConfig) error {
	constraint, err := cfg.Extension.GetShopwareVersionConstraint()

	if err != nil {
		return err
	}

	versions, err := getShopwareVersions()

	if err != nil {
		return err
	}

	vs := make([]*version.Version, 0)

	for _, r := range versions {
		v, err := version.NewVersion(r)
		if err != nil {
			continue
		}

		vs = append(vs, v)
	}

	sort.Sort(version.Collection(vs))

	matchingVersions := make([]*version.Version, 0)

	for _, v := range vs {
		if constraint.Check(v) {
			matchingVersions = append(matchingVersions, v)
		}
	}

	if len(matchingVersions) == 0 {
		matchingVersions = append(matchingVersions, version.Must(version.NewVersion("6.7.0.0")))
	}

	cfg.MinShopwareVersion = matchingVersions[0].String()
	cfg.MaxShopwareVersion = matchingVersions[len(matchingVersions)-1].String()

	return nil
}

type packagistResponse struct {
	Packages struct {
		Core []struct {
			Version string `json:"version_normalized"`
		} `json:"shopware/core"`
	} `json:"packages"`
}

func getShopwareVersions() ([]string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://repo.packagist.org/p2/shopware/core.json", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create composer version request: %w", err)
	}

	req.Header.Set("User-Agent", "Shopware Extension Verifier")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch composer versions: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("fetch composer versions: %s", resp.Status)
	}

	var pckResponse packagistResponse

	var versions []string

	if err := json.NewDecoder(resp.Body).Decode(&pckResponse); err != nil {
		return nil, fmt.Errorf("decode composer versions: %w", err)
	}

	for _, v := range pckResponse.Packages.Core {
		versions = append(versions, v.Version)
	}
	return versions, nil
}

func getAdminFolders(ext extension.Extension) []string {
	paths := []string{
		path.Join(ext.GetResourcesDir(), "app", "administration"),
	}

	for _, bundle := range ext.GetExtensionConfig().Build.ExtraBundles {
		paths = append(paths, path.Join(ext.GetRootDir(), bundle.Path, "Resources", "app", "administration"))
	}

	return filterNotExistingPaths(paths)
}

func getStorefrontFolders(ext extension.Extension) []string {
	paths := []string{
		path.Join(ext.GetResourcesDir(), "app", "storefront"),
	}

	for _, bundle := range ext.GetExtensionConfig().Build.ExtraBundles {
		paths = append(paths, path.Join(ext.GetRootDir(), bundle.Path, "Resources", "app", "storefront"))
	}

	return filterNotExistingPaths(paths)
}

func filterNotExistingPaths(paths []string) []string {
	filteredPaths := make([]string, 0)
	for _, p := range paths {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			filteredPaths = append(filteredPaths, p)
		}
	}

	return filteredPaths
}
