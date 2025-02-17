package tool

import (
	"context"
	"os"
	"os/exec"
	"path"
)

var ignoredPaths = `
package-lock.json
Resources/public/**
dist/**
Resources/store/**
`

type Prettier struct{}

func (b Prettier) Check(ctx context.Context, check *Check, config ToolConfig) error {
	return nil
}

func (b Prettier) Fix(ctx context.Context, config ToolConfig) error {
	return nil
}

func (b Prettier) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	rootDir := config.Extension.GetRootDir()

	if !path.IsAbs(rootDir) {
		rootDir = path.Join(cwd, rootDir)
	}

	if err := os.WriteFile(path.Join(rootDir, ".prettierignore"), []byte(ignoredPaths), 0644); err != nil {
		return err
	}

	args := []string{
		path.Join(cwd, "tools", "js", "node_modules", ".bin", "prettier"),
		"--config",
		path.Join(cwd, "tools", "js", ".prettierrc.js"),
		".",
	}

	if !dryRun {
		args = append(args, "--write")
	}

	cmd := exec.CommandContext(ctx, "node", args...)
	cmd.Dir = rootDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Remove(path.Join(rootDir, ".prettierignore"))
}

func init() {
	AddTool(Prettier{})
}
