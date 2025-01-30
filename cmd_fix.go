package main

import (
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var fixCommand = &cobra.Command{
	Use:   "fix [path]",
	Args:  cobra.ExactArgs(1),
	Short: "Fixes known issues in a Shopware extension",
	RunE: func(cmd *cobra.Command, args []string) error {
		toolCfg, err := guessExtension(args[0])

		if err != nil {
			return err
		}

		var gr errgroup.Group

		for _, tool := range availableTools {
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

func init() {
	rootCmd.AddCommand(fixCommand)
}
