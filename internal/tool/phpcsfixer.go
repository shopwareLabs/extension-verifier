package tool

import (
	"context"
	"os"
	"os/exec"
	"path"
)

type PHPCSFixer struct{}

func (p PHPCSFixer) Check(ctx context.Context, check *Check, config ToolConfig) error {
	return nil
}

func (p PHPCSFixer) Fix(ctx context.Context, config ToolConfig) error {
	return nil
}

func (p PHPCSFixer) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	if config.Extension.GetType() == "app" {
		return nil
	}

	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	rootDir := config.Extension.GetRootDir()

	if !path.IsAbs(rootDir) {
		rootDir = path.Join(cwd, rootDir)
	}

	args := []string{"fix", "--config", path.Join(cwd, "tools", "php-cs-fixer", ".php-cs-fixer.dist.php"), rootDir}
	if dryRun {
		args = append(args, "--dry-run")
	}
	cmd := exec.CommandContext(ctx, path.Join(cwd, "tools", "php-cs-fixer", "vendor", "bin", "php-cs-fixer"), args...)
	cmd.Dir = config.Extension.GetPath()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func init() {
	AddTool(PHPCSFixer{})
}
