package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
)

func runFixerOnString(fixer AdminTwigFixer, content string) (string, error) {
	nodes, err := html.NewParser(content)
	if err != nil {
		return "", err
	}

	err = fixer.Fix(nodes)
	if err != nil {
		return "", err
	}

	var buf strings.Builder

	for _, node := range nodes {
		buf.WriteString(node.Dump(0))
	}

	return buf.String(), nil
}
