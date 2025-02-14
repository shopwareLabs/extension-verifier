package main

import (
	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/shopware/shopware-cli/extension"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var formatCommand = &cobra.Command{
	Use:   "format [path]",
	Args:  cobra.ExactArgs(1),
	Short: "Formats the Shopware extension",
	RunE: func(cmd *cobra.Command, args []string) error {
		ext, err := extension.GetExtensionByFolder(args[0])

		if err != nil {
			return err
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

		toolCfg, err := tool.ConvertExtensionToToolConfig(ext)

		if err != nil {
			return err
		}

		var gr errgroup.Group

		for _, tool := range tool.GetTools() {
			gr.Go(func() error {
				return tool.Format(cmd.Context(), *toolCfg, dryRun)
			})
		}

		if err := gr.Wait(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(formatCommand)
	formatCommand.PersistentFlags().Bool("dry-run", false, "Dry run the formatting")
}
