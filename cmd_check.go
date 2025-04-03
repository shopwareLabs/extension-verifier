package main

import (
	"fmt"
	"os"

	"github.com/shopware/extension-verifier/internal/tool"
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

		if reportingFormat == "" {
			reportingFormat = detectDefaultReporter()
		}

		if err != nil {
			return err
		}

		stat, err := os.Stat(args[0])

		if err != nil {
			return err
		}

		var toolCfg *tool.ToolConfig

		if stat.IsDir() {
			if err := copyFiles(args[0], tmpDir); err != nil {
				return err
			}

			toolCfg, err = getToolConfig(tmpDir)
			if err != nil {
				return err
			}
		} else {
			ext, err := extension.GetExtensionByZip(args[0])

			if err != nil {
				return err
			}

			toolCfg, err = tool.ConvertExtensionToToolConfig(ext)

			if err != nil {
				return err
			}
		}

		toolCfg.CheckAgainst = checkAgainst

		defer func() {
			if err := os.RemoveAll(tmpDir); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to remove temporary directory: %v\n", err)
			}
		}()

		result := tool.NewCheck()

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
				return tool.Check(cmd.Context(), result, *toolCfg)
			})
		}

		if err := gr.Wait(); err != nil {
			return err
		}

		return doCheckReport(result.RemoveByIdentifier(toolCfg.ValidationIgnores), reportingFormat)
	},
}

func init() {
	rootCmd.AddCommand(checkCommand)
	checkCommand.PersistentFlags().String("reporter", "", "Reporting format (summary, json, github, junit, markdown)")
	checkCommand.PersistentFlags().String("check-against", "highest", "Check against Shopware Version (highest, lowest)")
	checkCommand.PersistentFlags().String("only", "", "Run only specific tools by name (comma-separated, e.g. phpstan,eslint)")
	checkCommand.PreRunE = func(cmd *cobra.Command, args []string) error {
		reporter, _ := cmd.Flags().GetString("reporter")
		if reporter != "summary" && reporter != "json" && reporter != "github" && reporter != "junit" && reporter != "markdown" && reporter != "" {
			return fmt.Errorf("invalid reporter format: %s. Must be either 'summary', 'json', 'github', 'junit' or 'markdown'", reporter)
		}

		mode, _ := cmd.Flags().GetString("check-against")
		if mode != "highest" && mode != "lowest" {
			return fmt.Errorf("invalid mode: %s. Must be either 'highest' or 'lowest'", mode)
		}

		return nil
	}
}
