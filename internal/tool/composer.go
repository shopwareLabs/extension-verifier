package tool

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/shopware/shopware-cli/extension"
)

func installComposerDeps(ext extension.Extension, checkAgainst string) error {
	rootDir := ext.GetPath()
	suggets := getComposerSuggets(ext)

	if _, err := os.Stat(path.Join(rootDir, "vendor")); os.IsNotExist(err) {
		if len(suggets) > 0 {
			additionalParams := []string{"require", "--prefer-dist", "--no-interaction", "--no-progress", "--no-plugins", "--no-scripts", "--ignore-platform-reqs"}
			for _, suggest := range suggets {
				additionalParams = append(additionalParams, fmt.Sprintf("%s:*", suggest))
			}

			composerInstall := exec.Command("composer", additionalParams...)
			composerInstall.Dir = rootDir

			log, err := composerInstall.CombinedOutput()

			if err != nil {
				os.Stderr.Write(log)
				return err
			}
		}

		additionalParams := []string{"update", "--prefer-dist", "--no-interaction", "--no-progress", "--no-plugins", "--no-scripts", "--ignore-platform-reqs"}

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

func getComposerSuggets(ext extension.Extension) []string {
	if inner, ok := ext.(*extension.PlatformPlugin); ok {
		suggests := make([]string, 0, len(inner.Composer.Suggest))
		for k := range inner.Composer.Suggest {
			suggests = append(suggests, k)
		}
		return suggests
	}

	if inner, ok := ext.(*extension.ShopwareBundle); ok {
		suggests := make([]string, 0, len(inner.Composer.Suggest))
		for k := range inner.Composer.Suggest {
			suggests = append(suggests, k)
		}
		return suggests
	}

	return []string{}
}
