package main

import (
	"encoding/json"
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

		if err != nil {
			return err
		}

		toolCfg, err := tool.ConvertExtensionToToolConfig(ext)

		if err != nil {
			return err
		}

		toolCfg.CheckAgainst = checkAgainst

		defer os.RemoveAll(tmpDir)

		result := tool.NewCheck()

		var gr errgroup.Group

		for _, tool := range tool.GetTools() {
			tool := tool
			gr.Go(func() error {
				return tool.Check(cmd.Context(), result, *toolCfg)
			})
		}

		if err := gr.Wait(); err != nil {
			return err
		}

		stepSummary := os.Getenv("GITHUB_STEP_SUMMARY")

		if stepSummary != "" {
			_ = os.WriteFile(stepSummary, []byte(convertResultsToMarkdown(result.Results)), 0644)
		}

		if reportingFormat == "summary" {
			// Group results by file
			fileGroups := make(map[string][]tool.CheckResult)
			for _, r := range result.Results {
				if r.Path == "" {
					r.Path = "general"
				}

				fileGroups[r.Path] = append(fileGroups[r.Path], r)
			}

			// Print results grouped by file
			totalProblems := 0
			errorCount := 0
			warningCount := 0

			for file, results := range fileGroups {
				fmt.Printf("\n%s\n", file)
				for _, r := range results {
					totalProblems++
					if r.Severity == "error" {
						errorCount++
					} else if r.Severity == "warning" {
						warningCount++
					}
					fmt.Printf("  %d  %-7s  %s  %s\n", r.Line, r.Severity, r.Message, r.Identifier)
				}
			}

			fmt.Printf("\n✖ %d problems (%d errors, %d warnings)\n", totalProblems, errorCount, warningCount)
		} else {
			j, err := json.Marshal(result)

			if err != nil {
				return err
			}

			os.Stdout.Write(j)
		}

		if result.HasErrors() {
			os.Exit(1)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(checkCommand)
	checkCommand.PersistentFlags().String("reporter", "summary", "Reporting format (summary, json)")
	checkCommand.PersistentFlags().String("check-against", "highest", "Check against Shopware Version (highest, lowest)")
	checkCommand.PreRunE = func(cmd *cobra.Command, args []string) error {
		reporter, _ := cmd.Flags().GetString("reporter")
		if reporter != "summary" && reporter != "json" {
			return fmt.Errorf("invalid reporter format: %s. Must be either 'summary' or 'json'", reporter)
		}

		mode, _ := cmd.Flags().GetString("check-against")
		if mode != "highest" && mode != "lowest" {
			return fmt.Errorf("invalid mode: %s. Must be either 'highest' or 'lowest'", mode)
		}

		return nil
	}
}
