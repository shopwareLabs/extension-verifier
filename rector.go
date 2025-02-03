package main

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
	cwd, err := os.Getwd()

	if err != nil {
		return err
	}

	rectorConfigFile := path.Join(cwd, "tools", "rector", "vendor", "frosh", "shopware-rector", "config", fmt.Sprintf("shopware-%s.0.php", config.MinShopwareVersion[0:3]))

	if err := installComposerDeps(config.Extension.GetPath(), "highest"); err != nil {
		return err
	}

	rector := exec.CommandContext(ctx, "php", "-dmemory_limit=2G", path.Join(cwd, "tools", "rector", "vendor", "bin", "rector"), "process", "--config", rectorConfigFile, "--autoload-file", path.Join("vendor", "autoload.php"), "src")
	rector.Dir = config.Extension.GetPath()

	log, _ := rector.CombinedOutput()

	fmt.Println(string(log))

	return nil
}
