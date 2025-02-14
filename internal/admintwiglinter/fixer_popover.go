package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type PopoverFixer struct{}

func init() {
	AddFixer(PopoverFixer{})
}

func (p PopoverFixer) Check(node []html.Node) []CheckError {
	var checkErrors []CheckError

	html.TraverseNode(node, func(node *html.ElementNode) {
		if node.Tag == "sw-popover" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-popover is deprecated, use mt-floating-ui instead",
				Severity:   "error",
				Identifier: "sw-popover",
			})
		}
	})

	return checkErrors
}

func (p PopoverFixer) Supports(version *version.Version) bool {
	return shopware67Constraint.Check(version)
}

func (p PopoverFixer) Fix(node []html.Node) error {
	html.TraverseNode(node, func(node *html.ElementNode) {
		if node.Tag == "sw-popover" {
			node.Tag = "mt-floating-ui"

			hasVIf := false

			for n, attr := range node.Attributes {
				if attr.Key == "v-if" {
					node.Attributes[n].Key = ":isOpened"
					hasVIf = true
				}

				if attr.Key == ":zIndex" || attr.Key == ":resizeWidth" {
					node.Attributes = append(node.Attributes[:n], node.Attributes[n+1:]...)
				}
			}

			if !hasVIf {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   ":isOpened",
					Value: "true",
				})
			}
		}
	})

	return nil
}
