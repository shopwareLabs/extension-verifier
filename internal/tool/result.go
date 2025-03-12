package tool

import (
	"sync"
)

type Check struct {
	Results []CheckResult `json:"results"`
	mutex   sync.Mutex
}

func NewCheck() *Check {
	return &Check{
		Results: []CheckResult{},
	}
}

func (c *Check) AddResult(result CheckResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Results = append(c.Results, result)
}

func (c *Check) HasErrors() bool {
	for _, r := range c.Results {
		if r.Severity == "error" {
			return true
		}
	}

	return false
}

func (c *Check) RemoveByIdentifier(ignores []ToolConfigIgnore) *Check {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var filtered []CheckResult
	for _, r := range c.Results {
		shouldKeep := true
		for _, ignore := range ignores {
			if r.Identifier == ignore.Identifier && (r.Path == ignore.Path || ignore.Path == "") {
				shouldKeep = false
				break
			}
		}
		if shouldKeep {
			filtered = append(filtered, r)
		}
	}
	c.Results = filtered

	return c
}

type CheckResult struct {
	// The path to the file that was checked
	Path string `json:"path"`
	// The line number of the issue
	Line    int    `json:"line"`
	Message string `json:"message"`
	// The severity of the issue
	Severity string `json:"severity"`

	Identifier string `json:"identifier"`
}
