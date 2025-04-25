package main

import (
	"context"
	"fmt"
	"os"

	"github.com/shopware/shopware-cli/logging"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sw-extension-verifier",
	Short: "A CLI tool to refactor and check the quality of Shopware extensions",
}

func main() {
	verbose := false

	if err := rootCmd.ParseFlags(os.Args); err == nil {
		verbose, _ = rootCmd.PersistentFlags().GetBool("verbose")
	}

	ctx := logging.WithLogger(context.Background(), logging.NewLogger(verbose))

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().Bool("verbose", false, "show debug output")
}
