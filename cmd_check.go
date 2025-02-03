package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss/table"
	"github.com/shopware/shopware-cli/extension"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var checkCommand = &cobra.Command{
	Use:   "check [path]",
	Args:  cobra.ExactArgs(1),
	Short: "Check the quality of a Shopware extension",
	RunE: func(cmd *cobra.Command, args []string) error {
		reportingFormat, _ := cmd.Flags().GetString("reporter")
		checkAgainst, _ := cmd.Flags().GetString("check-against")
		tmpDir, err := os.MkdirTemp(os.TempDir(), "analyse-extension-*")

		if err != nil {
			return err
		}

		stat, err := os.Stat(args[0])

		if err != nil {
			return err
		}

		var ext extension.Extension

		if stat.IsDir() {
			if err := copyFiles(args[0], tmpDir); err != nil {
				return err
			}

			ext, err = extension.GetExtensionByFolder(tmpDir)
		} else {
			ext, err = extension.GetExtensionByZip(args[0])
		}

		toolCfg, err := convertExtensionToToolConfig(ext)

		if err != nil {
			return err
		}

		toolCfg.CheckAgainst = checkAgainst

		defer os.RemoveAll(tmpDir)

		result := newCheck()

		var gr errgroup.Group

		for _, tool := range availableTools {
			tool := tool
			gr.Go(func() error {
				return tool.Check(cmd.Context(), result, *toolCfg)
			})
		}

		if err := gr.Wait(); err != nil {
			return err
		}

		if reportingFormat == "table" {
			t := table.New().Headers("Severity", "Identifier", "File", "Message")

			for _, r := range result.Results {
				t.Row(r.Severity, r.Identifier, fmt.Sprintf("%s:%d", r.Path, r.Line), r.Message)
			}

			fmt.Println(t.String())
		} else {
			j, err := json.Marshal(result)

			if err != nil {
				return err
			}

			os.Stdout.Write(j)
		}

		if len(result.Results) > 0 {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCommand)
	checkCommand.PersistentFlags().String("reporter", "table", "Reporting format (table, json)")
	checkCommand.PersistentFlags().String("check-against", "highest", "Check against Shopware Version (highest, lowest)")
	checkCommand.PreRunE = func(cmd *cobra.Command, args []string) error {
		reporter, _ := cmd.Flags().GetString("reporter")
		if reporter != "table" && reporter != "json" {
			return fmt.Errorf("invalid reporter format: %s. Must be either 'table' or 'json'", reporter)
		}

		mode, _ := cmd.Flags().GetString("check-against")
		if mode != "highest" && mode != "lowest" {
			return fmt.Errorf("invalid mode: %s. Must be either 'highest' or 'lowest'", mode)
		}

		return nil
	}
}
