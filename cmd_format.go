package main

import (
	"github.com/shopware/extension-verifier/internal/tool"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var formatCommand = &cobra.Command{
	Use:   "format [path]",
	Args:  cobra.ExactArgs(1),
	Short: "Formats the Shopware extension",
	RunE: func(cmd *cobra.Command, args []string) error {
		toolCfg, err := getToolConfig(args[0])
		if err != nil {
			return err
		}

		dryRun, _ := cmd.Flags().GetBool("dry-run")

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
	formatCommand.PersistentFlags().String("only", "", "Run only specific tools by name (comma-separated, e.g. phpstan,eslint)")
}
