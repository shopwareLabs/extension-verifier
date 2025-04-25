package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/shopware-cli/logging"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	allowNonGit bool
	fixCommand  = &cobra.Command{
		Use:   "fix [path]",
		Args:  cobra.ExactArgs(1),
		Short: "Fixes known issues in a Shopware extension",
		RunE: func(cmd *cobra.Command, args []string) error {
			gitPath := filepath.Join(args[0], ".git")
			if !allowNonGit {
				if stat, err := os.Stat(gitPath); err != nil || !stat.IsDir() {
					return fmt.Errorf("provided folder is not a git repository. Use --allow-non-git flag to run anyway")
				}
			}

			toolCfg, err := getToolConfig(args[0])
			if err != nil {
				return err
			}

			logging.FromContext(cmd.Context()).Debugf("Running fixes for Shopware version: %s", toolCfg.MinShopwareVersion)

			var gr errgroup.Group

			tools := tool.GetTools()
			only, _ := cmd.Flags().GetString("only")

			tools, err = filterTools(tools, only)
			if err != nil {
				return err
			}

			for _, tool := range tools {
				tool := tool
				gr.Go(func() error {
					return tool.Fix(cmd.Context(), *toolCfg)
				})
			}

			if err := gr.Wait(); err != nil {
				return err
			}

			return nil
		},
	}
)

func init() {
	fixCommand.Flags().BoolVar(&allowNonGit, "allow-non-git", false, "Allow running the fix command on non-git repositories")
	fixCommand.Flags().String("only", "", "Run only specific tools by name (comma-separated, e.g. phpstan,eslint)")
	rootCmd.AddCommand(fixCommand)
}
