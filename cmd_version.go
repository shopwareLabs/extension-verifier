package main

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var version = "dev"

func init() {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("extension-verifier version %s\n", version)
			fmt.Printf("  Go version: %s\n", runtime.Version())
			fmt.Printf("  OS/Arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}

	rootCmd.AddCommand(versionCmd)
}