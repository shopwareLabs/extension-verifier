package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type SwitchFixer struct{}

func init() {
	AddFixer(SwitchFixer{})
}

func (s SwitchFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-switch-field" {
			errs = append(errs, CheckError{
				Message:    "sw-switch-field is removed, use mt-switch instead. Review conversion for props, events and slots.",
				Severity:   "error",
				Identifier: "sw-switch-field",
			})
		}
	})
	return errs
}

func (s SwitchFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (s SwitchFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-switch-field" {
			node.Tag = "mt-switch"
			var newAttrs []html.Attribute
			// Process attribute conversions.
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "noMarginTop":
					newAttrs = append(newAttrs, html.Attribute{Key: "removeTopMargin"})
				case "size", "id", "ghostValue", "padded", "partlyChecked":
					// remove these attributes
				case "value":
					newAttrs = append(newAttrs, html.Attribute{Key: "checked", Value: attr.Value})
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process children for slot conversion.
			var labelText string
			var otherChildren []html.Node
			for _, child := range node.Children {
				// Check if child is a slot element.
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					for _, a := range elem.Attributes {
						if a.Key == "#label" {
							var sb strings.Builder
							for _, inner := range elem.Children {
								sb.WriteString(strings.TrimSpace(inner.Dump()))
							}
							labelText = sb.String()
							goto NextChild
						}
						if a.Key == "#hint" {
							goto NextChild
						}
					}
				}
				otherChildren = append(otherChildren, child)
			NextChild:
			}
			// Remove all slot children.
			node.Children = otherChildren
			// If label slot found, add label attribute.
			if labelText != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: labelText,
				})
			}
		}
	})
	return nil
}
