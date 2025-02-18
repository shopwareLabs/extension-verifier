package tool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type Rector struct{}

func (r Rector) Check(ctx context.Context, check *Check, config ToolConfig) error {
	return nil
}

func (r Rector) Fix(ctx context.Context, config ToolConfig) error {
	if config.Extension.GetType() == "app" {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Check and remove existing vendor/composer.lock
	vendorPath := path.Join(config.Extension.GetPath(), "vendor")
	composerLockPath := path.Join(config.Extension.GetPath(), "composer.lock")

	if _, err := os.Stat(vendorPath); err == nil {
		if err := os.RemoveAll(vendorPath); err != nil {
			return err
		}
	}
	if _, err := os.Stat(composerLockPath); err == nil {
		if err := os.Remove(composerLockPath); err != nil {
			return err
		}
	}

	rectorConfigFile := path.Join(cwd, "tools", "php", "vendor", "frosh", "shopware-rector", "config", fmt.Sprintf("shopware-%s.0.php", config.MinShopwareVersion[0:3]))

	if err := installComposerDeps(config.Extension.GetPath(), "highest"); err != nil {
		return err
	}

	rector := exec.CommandContext(ctx, "php", "-dmemory_limit=2G", path.Join(cwd, "tools", "php", "vendor", "bin", "rector"), "process", "--config", rectorConfigFile, "--autoload-file", path.Join("vendor", "autoload.php"), "src")
	rector.Dir = config.Extension.GetPath()

	log, _ := rector.CombinedOutput()
	fmt.Println(string(log))

	// Cleanup after execution
	if err := os.RemoveAll(vendorPath); err != nil {
		return err
	}
	if err := os.Remove(composerLockPath); err != nil {
		return err
	}

	return nil
}

func (r Rector) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}

func init() {
	AddTool(Rector{})
}
