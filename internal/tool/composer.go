package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
)

func installComposerDeps(rootDir string, checkAgainst string) error {
	suggets := getComposerSuggets(rootDir)

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

func getComposerSuggets(rootDir string) []string {
	composerJSON, err := os.ReadFile(path.Join(rootDir, "composer.json"))
	if err != nil {
		return []string{}
	}

	var composerJSONData map[string]interface{}
	if err := json.Unmarshal(composerJSON, &composerJSONData); err != nil {
		return []string{}
	}

	if composerJSONData["suggest"] == nil {
		return []string{}
	}

	suggests := make([]string, 0, len(composerJSONData["suggest"].(map[string]interface{})))
	for k := range composerJSONData["suggest"].(map[string]interface{}) {
		suggests = append(suggests, k)
	}

	return suggests
}
