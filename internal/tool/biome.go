package tool

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path"
)

//go:embed biome.json
var biomeConfig []byte

type Biome struct{}

func (b Biome) Check(ctx context.Context, check *Check, config ToolConfig) error {
	return nil
}

func (b Biome) Fix(ctx context.Context, config ToolConfig) error {
	return nil
}

func (b Biome) Format(ctx context.Context, config ToolConfig, dryRun bool) error {
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	if err := os.WriteFile(path.Join(cwd, "biome.json"), biomeConfig, 0644); err != nil {
		return err
	}

	rootDir := config.Extension.GetRootDir()

	if !path.IsAbs(rootDir) {
		rootDir = path.Join(cwd, rootDir)
	}

	args := []string{"format", fmt.Sprintf("--config-path=%s", path.Join(cwd, "biome.json"))}

	if !dryRun {
		args = append(args, "--write")
	}

	cmd := exec.CommandContext(ctx, "biome", args...)
	cmd.Dir = rootDir
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Remove(path.Join(cwd, "biome.json"))
}

func init() {
	AddTool(Biome{})
}
