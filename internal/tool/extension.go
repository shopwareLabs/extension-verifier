package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"

	"github.com/shopware/shopware-cli/extension"
	"github.com/shopware/shopware-cli/version"
)

func ConvertExtensionToToolConfig(ext extension.Extension) (*ToolConfig, error) {
	cfg := &ToolConfig{
		Extension:         ext,
		ValidationIgnores: ext.GetExtensionConfig().Validation.Ignore,
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
