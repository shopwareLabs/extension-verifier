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
	// Apps don't have an composer.json file, skip them
	if _, err := os.Stat(path.Join(config.RootDir, "composer.json")); err != nil {
		return nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Backup composer.json
	composerJSONPath := path.Join(config.RootDir, "composer.json")
	var backupData []byte
	if _, err := os.Stat(composerJSONPath); err == nil {
		backupData, err = os.ReadFile(composerJSONPath)
		if err != nil {
			return fmt.Errorf("failed to backup composer.json: %w", err)
		}
	}

	// Check and remove existing vendor/composer.lock
	vendorPath := path.Join(config.RootDir, "vendor")
	composerLockPath := path.Join(config.RootDir, "composer.lock")

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

	if err := installComposerDeps(config.RootDir, "highest"); err != nil {
		return err
	}

	for _, sourceDirectory := range config.SourceDirectories {
		rector := exec.CommandContext(ctx, "php", "-dmemory_limit=2G", path.Join(cwd, "tools", "php", "vendor", "bin", "rector"), "process", "--config", rectorConfigFile, "--autoload-file", path.Join("vendor", "autoload.php"), sourceDirectory)
		rector.Dir = config.RootDir

		log, _ := rector.CombinedOutput()
		fmt.Println(string(log))
	}

	// Cleanup after execution
	if err := os.RemoveAll(vendorPath); err != nil {
		return err
	}
	if err := os.Remove(composerLockPath); err != nil {
		return err
	}

	// Restore composer.json
	if backupData != nil {
		if err := os.WriteFile(composerJSONPath, backupData, 0644); err != nil {
			return fmt.Errorf("failed to restore composer.json: %w", err)
		}
	}

	return nil
}

func (r Rector) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	return nil
}

func init() {
	AddTool(Rector{})
}
