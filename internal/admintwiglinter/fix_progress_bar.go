package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type ProgressBarFixer struct{}

func init() {
	AddFixer(ProgressBarFixer{})
}

func (p ProgressBarFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-progress-bar" {
			errors = append(errors, CheckError{
				Message:    "sw-progress-bar is removed, use mt-progress-bar instead.",
				Severity:   "error",
				Identifier: "sw-progress-bar",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (p ProgressBarFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (p ProgressBarFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-progress-bar" {
			node.Tag = "mt-progress-bar"
			// Update attribute keys.
			for i, attr := range node.Attributes {
				switch attr.Key {
				case "value":
					node.Attributes[i].Key = "modelValue"
				case "v-model:value":
					node.Attributes[i].Key = "v-model"
				case "@update:value":
					node.Attributes[i].Key = "@update:modelValue"
				}
			}
		}
	})
	return nil
}
