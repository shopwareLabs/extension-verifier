package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/shopware-cli/extension"
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

			ext, err := extension.GetExtensionByFolder(args[0])

			if err != nil {
				return err
			}

			toolCfg, err := tool.ConvertExtensionToToolConfig(ext)

			if err != nil {
				return err
			}

			var gr errgroup.Group

			for _, tool := range tool.GetTools() {
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
	rootCmd.AddCommand(fixCommand)
}
