package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type UrlFieldFixer struct{}

func init() {
	AddFixer(UrlFieldFixer{})
}

func (u UrlFieldFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-url-field" {
			errors = append(errors, CheckError{
				Message:    "sw-url-field is removed, use mt-url-field instead. Review conversion for props, events, label and hint slot.",
				Severity:   "error",
				Identifier: "sw-url-field",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (u UrlFieldFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (u UrlFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-url-field" {
			node.Tag = "mt-url-field"
			var newAttrs []html.Attribute
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "value":
					attr.Key = "model-value"
					newAttrs = append(newAttrs, attr)
				case "v-model:value":
					attr.Key = "v-model"
					newAttrs = append(newAttrs, attr)
				case "@update:value":
					attr.Key = "@update:model-value"
					newAttrs = append(newAttrs, attr)
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process slot conversion.
			label := ""
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					for _, a := range elem.Attributes {
						if a.Key == "#label" {
							var content strings.Builder
							for _, inner := range elem.Children {
								content.WriteString(strings.TrimSpace(inner.Dump()))
							}
							label = content.String()
						}
						if a.Key == "#hint" {
							// Skip hint slot.
							continue
						}
					}
				}
			}
			// Remove all children; label was processed, and hint slot is dropped.
			node.Children = []html.Node{}
			if label != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: label,
				})
			}
		}
	})
	return nil
}
