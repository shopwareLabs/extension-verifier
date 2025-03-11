package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
)

var shopware67Constraint = version.MustConstraints(version.NewConstraint(">=6.7.0"))

var availableFixers = []AdminTwigFixer{}

func AddFixer(fixer AdminTwigFixer) {
	availableFixers = append(availableFixers, fixer)
}

type CheckError struct {
	Message    string
	Severity   string
	Identifier string
	Line       int
}

func GetFixers(version *version.Version) []AdminTwigFixer {
	fixers := []AdminTwigFixer{}
	for _, fixer := range availableFixers {
		if fixer.Supports(version) {
			fixers = append(fixers, fixer)
		}
	}

	return fixers
}

type AdminTwigFixer interface {
	Check(node []html.Node) []CheckError
	Supports(version *version.Version) bool
	Fix(node []html.Node) error
}
