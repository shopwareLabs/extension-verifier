package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"

	"github.com/shopware/extension-verifier/internal/tool"
)

func detectDefaultReporter() string {
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		return "github"
	}

	return "summary"
}

func doCheckReport(result *tool.Check, reportingFormat string) error {
	switch reportingFormat {
	case "summary":
		return doSummaryReport(result)
	case "json":
		return doJSONReport(result)
	case "github":
		return doGitHubReport(result)
	case "markdown":
		return doMarkdownReport(result)
	case "junit":
		return doJUnitReport(result)
	}

	if result.HasErrors() {
		os.Exit(1)
	}

	return nil
}

func doSummaryReport(result *tool.Check) error {
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

	return nil
}

func doJSONReport(result *tool.Check) error {
	j, err := json.Marshal(result)

	if err != nil {
		return err
	}

	os.Stdout.Write(j)

	return nil
}

func doGitHubReport(result *tool.Check) error {
	stepSummary := os.Getenv("GITHUB_STEP_SUMMARY")

	if stepSummary != "" {
		if err := os.WriteFile(stepSummary, []byte(convertResultsToMarkdown(result.Results)), 0644); err != nil {
			return fmt.Errorf("failed to write step summary: %w", err)
		}
	}

	for _, res := range result.Results {
		if res.Line == 0 {
			fmt.Printf("::%s file=%s::%s\n", res.Severity, res.Path, res.Message)
		} else {
			fmt.Printf("::%s file=%s,line=%d::%s\n", res.Severity, res.Path, res.Line, res.Message)
		}
	}

	return nil
}

func doJUnitReport(result *tool.Check) error {
	type testcase struct {
		XMLName   xml.Name `xml:"testcase"`
		Name      string   `xml:"name,attr"`
		Classname string   `xml:"classname,attr"`
		Time      string   `xml:"time,attr"`
		Failure   *struct {
			Message string `xml:"message,attr"`
			Type    string `xml:"type,attr"`
			Content string `xml:",chardata"`
		} `xml:"failure,omitempty"`
	}

	type testsuite struct {
		XMLName   xml.Name   `xml:"testsuite"`
		Name      string     `xml:"name,attr"`
		Tests     int        `xml:"tests,attr"`
		Failures  int        `xml:"failures,attr"`
		Errors    int        `xml:"errors,attr"`
		Time      string     `xml:"time,attr"`
		Testcases []testcase `xml:"testcase"`
	}

	type testsuites struct {
		XMLName    xml.Name    `xml:"testsuites"`
		Testsuites []testsuite `xml:"testsuite"`
	}

	// Create a test case for each result
	testcases := make([]testcase, 0, len(result.Results))
	failures := 0

	for _, res := range result.Results {
		tc := testcase{
			Name:      res.Identifier,
			Classname: res.Path,
			Time:      "0.000", // No timing information available
		}

		// Add failure information if severity is not "notice"
		if res.Severity != "notice" {
			failures++
			tc.Failure = &struct {
				Message string `xml:"message,attr"`
				Type    string `xml:"type,attr"`
				Content string `xml:",chardata"`
			}{
				Message: res.Message,
				Type:    res.Severity,
				Content: fmt.Sprintf("Line: %d\nMessage: %s", res.Line, res.Message),
			}
		}

		testcases = append(testcases, tc)
	}

	// Create the test suite
	ts := testsuite{
		Name:      "Extension Verification",
		Tests:     len(testcases),
		Failures:  failures,
		Errors:    0,       // We don't distinguish between failures and errors
		Time:      "0.000", // No timing information available
		Testcases: testcases,
	}

	// Create the root element
	root := testsuites{
		Testsuites: []testsuite{ts},
	}

	// Marshal to XML
	output, err := xml.MarshalIndent(root, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JUnit XML: %w", err)
	}

	// Add XML header
	output = append([]byte(xml.Header), output...)

	// Write to stdout
	_, err = os.Stdout.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write JUnit XML: %w", err)
	}

	return nil
}

func doMarkdownReport(result *tool.Check) error {
	os.Stdout.Write([]byte(convertResultsToMarkdown(result.Results)))

	return nil
}

func convertResultsToMarkdown(check []tool.CheckResult) string {
	var builder strings.Builder

	builder.WriteString("# Results\n\n")

	builder.WriteString("| Severity | Identifier | File | Message | \n")
	builder.WriteString("| --- | --- | --- | --- |\n")

	for _, result := range check {
		builder.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", result.Severity, result.Identifier, result.Path, result.Message))
	}

	builder.WriteString("\n")

	return builder.String()
}
