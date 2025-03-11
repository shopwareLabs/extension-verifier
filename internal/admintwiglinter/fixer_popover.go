package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
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
				Line:       node.Line,
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
			var newAttrs html.NodeList

			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case "v-if":
						attr.Key = ":isOpened"
						newAttrs = append(newAttrs, attr)
						hasVIf = true
					case ":zIndex", ":resizeWidth":
						// Skip these attributes
					default:
						newAttrs = append(newAttrs, attr)
					}
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}

			if !hasVIf {
				newAttrs = append(newAttrs, html.Attribute{
					Key:   ":isOpened",
					Value: "true",
				})
			}

			node.Attributes = newAttrs
		}
	})

	return nil
}
