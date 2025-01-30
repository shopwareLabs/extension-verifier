package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/shopware/shopware-cli/version"
)

func guessExtension(rootDir string) (*ToolConfig, error) {
	if _, err := os.Stat(filepath.Join(rootDir, "composer.json")); err == nil {
		return guessByPlugin(rootDir)
	}

	return nil, fmt.Errorf("could not guess extension type")
}

type ComposerJson struct {
	Require map[string]string `json:"require"`
}

func guessByPlugin(rootDir string) (*ToolConfig, error) {
	var composerJson ComposerJson

	jsonFile, err := os.ReadFile(filepath.Join(rootDir, "composer.json"))

	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(jsonFile, &composerJson); err != nil {
		return nil, err
	}

	if _, ok := composerJson.Require["shopware/core"]; !ok {
		return nil, fmt.Errorf("shopware/core requirement is missing")
	}

	cfg := &ToolConfig{RootDir: rootDir, ShopwareVersionConstraint: composerJson.Require["shopware/core"]}

	if err := determineVersionRange(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func determineVersionRange(cfg *ToolConfig) error {
	constraint := version.MustConstraints(version.NewConstraint(cfg.ShopwareVersionConstraint))

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
		return fmt.Errorf("the composer constraint does not match any shopware release")
	}

	cfg.MinShopwareVersion = matchingVersions[0].String()
	cfg.MaxShopwareVersion = matchingVersions[len(matchingVersions)-1].String()

	return nil
}

func getShopwareVersions() ([]string, error) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, "https://shopware.github.io/shopware-cli/versions.json", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("create composer version request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetch composer versions: %w", err)
	}
	defer func() {
		resp.Body.Close()
	}()

	versionString, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read version body: %w", err)
	}

	var versions []string
	if err := json.Unmarshal(versionString, &versions); err != nil {
		return nil, fmt.Errorf("unmarshal composer versions: %w", err)
	}

	return versions, nil
}
