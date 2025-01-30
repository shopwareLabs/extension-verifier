package main

import (
	"os"
	"os/exec"
	"path"
)

func installComposerDeps(rootDir, checkAgainst string) error {
	if _, err := os.Stat(path.Join(rootDir, "vendor")); os.IsNotExist(err) {
		additionalParams := []string{"update", "--no-interaction", "--no-progress", "--no-plugins", "--no-scripts"}

		if checkAgainst == "lowest" {
			additionalParams = append(additionalParams, "--prefer-lowest")
		}

		composerInstall := exec.Command("composer", additionalParams...)
		composerInstall.Dir = rootDir

		log, err := composerInstall.CombinedOutput()

		if err != nil {
			os.Stderr.Write(log)
			return err
		}
	}

	return nil
}
