package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type AlertFixer struct{}

func init() {
	AddFixer(AlertFixer{})
}

func (a AlertFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-alert" {
			errors = append(errors, CheckError{
				Message:    "sw-alert is removed, use mt-banner instead. Please review conversion for variant changes.",
				Severity:   "error",
				Identifier: "sw-alert",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (a AlertFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (a AlertFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-alert" {
			node.Tag = "mt-banner"
			var newAttrs []html.Attribute

			for _, attr := range node.Attributes {
				if attr.Key == "variant" {
					switch attr.Value {
					case "success":
						attr.Value = "positive"
						newAttrs = append(newAttrs, attr)
					case "error":
						attr.Value = "critical"
						newAttrs = append(newAttrs, attr)
					case "warning":
						attr.Value = "attention"
						newAttrs = append(newAttrs, attr)
					case "info":
						// Keep info as is
						newAttrs = append(newAttrs, attr)
					default:
						// Keep any other variants unchanged
						newAttrs = append(newAttrs, attr)
					}
				} else {
					// Preserve all other attributes
					newAttrs = append(newAttrs, attr)
				}
			}

			node.Attributes = newAttrs
		}
	})
	return nil
}
