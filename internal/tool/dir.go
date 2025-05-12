package tool

import (
	"os"
	"path"
)

var toolDirectory = ""

func SetToolDirectory(dir string) {
	toolDirectory = dir
}

func GetToolDirectory() string {
	if toolDirectory != "" {
		return toolDirectory
	}

	cwd, err := os.Getwd()

	if err != nil {
		return ""
	}

	toolDirectory = path.Join(cwd, "tools")

	return toolDirectory
}
