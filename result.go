package main

import "sync"

type Check struct {
	Results []CheckResult `json:"results"`
	mutex   sync.Mutex
}

func newCheck() *Check {
	return &Check{
		Results: []CheckResult{},
	}
}

func (c *Check) AddResult(result CheckResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.Results = append(c.Results, result)
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
