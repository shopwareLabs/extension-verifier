package tool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCheck(t *testing.T) {
	check := NewCheck()
	assert.NotNil(t, check)
	assert.Empty(t, check.Results)
}

func TestAddResult(t *testing.T) {
	check := NewCheck()
	result := CheckResult{
		Path:       "test.go",
		Line:       1,
		Message:    "test message",
		Severity:   "error",
		Identifier: "TEST001",
	}

	check.AddResult(result)
	assert.Len(t, check.Results, 1)
	assert.Equal(t, result, check.Results[0])
}

func TestHasErrors(t *testing.T) {
	tests := []struct {
		name     string
		results  []CheckResult
		expected bool
	}{
		{
			name:     "no results",
			results:  []CheckResult{},
			expected: false,
		},
		{
			name: "no errors",
			results: []CheckResult{
				{Severity: "warning"},
				{Severity: "info"},
			},
			expected: false,
		},
		{
			name: "has errors",
			results: []CheckResult{
				{Severity: "warning"},
				{Severity: "error"},
				{Severity: "info"},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewCheck()
			for _, result := range tt.results {
				check.AddResult(result)
			}
			assert.Equal(t, tt.expected, check.HasErrors())
		})
	}
}

func TestRemoveByIdentifier(t *testing.T) {
	tests := []struct {
		name            string
		initialResults  []CheckResult
		ignores         []ToolConfigIgnore
		expectedResults []CheckResult
	}{
		{
			name: "remove single result",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001"},
				{Path: "file2.go", Identifier: "TEST002"},
			},
			ignores: []ToolConfigIgnore{
				{Path: "file1.go", Identifier: "TEST001"},
			},
			expectedResults: []CheckResult{
				{Path: "file2.go", Identifier: "TEST002"},
			},
		},
		{
			name: "remove by identifier only",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001"},
				{Path: "file2.go", Identifier: "TEST001"},
			},
			ignores: []ToolConfigIgnore{
				{Identifier: "TEST001"},
			},
			expectedResults: []CheckResult{},
		},
		{
			name: "no matches",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001"},
				{Path: "file2.go", Identifier: "TEST002"},
			},
			ignores: []ToolConfigIgnore{
				{Path: "file3.go", Identifier: "TEST003"},
			},
			expectedResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001"},
				{Path: "file2.go", Identifier: "TEST002"},
			},
		},
		{
			name: "identifier with message should not ignore all",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001", Message: "error 1"},
				{Path: "file2.go", Identifier: "TEST001", Message: "error 2"},
			},
			ignores: []ToolConfigIgnore{
				{Identifier: "TEST001", Message: "error 1"},
			},
			expectedResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001", Message: "error 1"},
				{Path: "file2.go", Identifier: "TEST001", Message: "error 2"},
			},
		},
		{
			name: "identifier only should ignore all matching errors",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001", Message: "error 1"},
				{Path: "file2.go", Identifier: "TEST001", Message: "error 2"},
				{Path: "file3.go", Identifier: "TEST002", Message: "error 3"},
			},
			ignores: []ToolConfigIgnore{
				{Identifier: "TEST001"},
			},
			expectedResults: []CheckResult{
				{Path: "file3.go", Identifier: "TEST002", Message: "error 3"},
			},
		},
		{
			name: "multiple ignores with different conditions",
			initialResults: []CheckResult{
				{Path: "file1.go", Identifier: "TEST001", Message: "error 1"},
				{Path: "file2.go", Identifier: "TEST001", Message: "error 2"},
				{Path: "file3.go", Identifier: "TEST002", Message: "error 3"},
				{Path: "file4.go", Identifier: "TEST003", Message: "error 4"},
			},
			ignores: []ToolConfigIgnore{
				{Identifier: "TEST001"},                   // Should remove all TEST001
				{Path: "file3.go", Identifier: "TEST002"}, // Should remove only TEST002 in file3.go
				{Message: "error 4"},                      // Should remove anything with this message
			},
			expectedResults: []CheckResult{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewCheck()
			for _, result := range tt.initialResults {
				check.AddResult(result)
			}

			check.RemoveByIdentifier(tt.ignores)
			assert.ElementsMatch(t, tt.expectedResults, check.Results)
		})
	}
}

func TestRemoveByMessage(t *testing.T) {
	tests := []struct {
		name            string
		initialResults  []CheckResult
		ignores         []ToolConfigIgnore
		expectedResults []CheckResult
	}{
		{
			name: "remove_single_result_by_exact_message",
			initialResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
			ignores: []ToolConfigIgnore{
				{Path: "file1.go", Message: "test error message"},
			},
			expectedResults: []CheckResult{
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
		},
		{
			name: "remove_by_message_only",
			initialResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "test error message", Severity: "error"},
			},
			ignores: []ToolConfigIgnore{
				{Message: "test error message"},
			},
			expectedResults: []CheckResult{},
		},
		{
			name: "remove_by_partial_message_match",
			initialResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
			ignores: []ToolConfigIgnore{
				{Message: "test error"},
			},
			expectedResults: []CheckResult{
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
		},
		{
			name: "no_matches",
			initialResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
			ignores: []ToolConfigIgnore{
				{Message: "non-existent message"},
			},
			expectedResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "another error message", Severity: "error"},
			},
		},
		{
			name: "remove_by_path_and_partial_message",
			initialResults: []CheckResult{
				{Path: "file1.go", Line: 1, Message: "test error message", Severity: "error"},
				{Path: "file2.go", Line: 2, Message: "test error message", Severity: "error"},
			},
			ignores: []ToolConfigIgnore{
				{Path: "file1.go", Message: "test"},
			},
			expectedResults: []CheckResult{
				{Path: "file2.go", Line: 2, Message: "test error message", Severity: "error"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewCheck()
			for _, result := range tt.initialResults {
				check.AddResult(result)
			}

			check.RemoveByIdentifier(tt.ignores)
			assert.ElementsMatch(t, tt.expectedResults, check.Results)
		})
	}
}
