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
		name           string
		initialResults []CheckResult
		ignores        []ToolConfigIgnore
		expectedCount  int
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
			expectedCount: 1,
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
			expectedCount: 0,
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
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			check := NewCheck()
			for _, result := range tt.initialResults {
				check.AddResult(result)
			}

			check.RemoveByIdentifier(tt.ignores)
			assert.Len(t, check.Results, tt.expectedCount)
		})
	}
}
