package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type CheckboxFieldFixer struct{}

func init() {
	AddFixer(CheckboxFieldFixer{})
}

func (c CheckboxFieldFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-checkbox-field" {
			errs = append(errs, CheckError{
				Message:    "sw-checkbox-field is removed, use mt-checkbox instead. Review conversion for props, events and slots.",
				Severity:   "error",
				Identifier: "sw-checkbox-field",
			})
		}
	})
	return errs
}

func (c CheckboxFieldFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (c CheckboxFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-checkbox-field" {
			node.Tag = "mt-checkbox"
			var newAttrs []html.Attribute
			// Process attribute conversions.
			for _, attr := range node.Attributes {
				switch attr.Key {
				case ":value":
					newAttrs = append(newAttrs, html.Attribute{Key: ":checked", Value: attr.Value})
				case "v-model":
					newAttrs = append(newAttrs, html.Attribute{Key: "v-model:checked", Value: attr.Value})
				case "id", "ghostValue", "padded":
					// remove these attributes without replacement
				case "partlyChecked":
					newAttrs = append(newAttrs, html.Attribute{Key: "partial"})
				case "@update:value":
					newAttrs = append(newAttrs, html.Attribute{Key: "@update:checked", Value: attr.Value})
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process children for slot conversion.
			var labelText string
			var remainingChildren []html.Node
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					// Handle label slot.
					for _, a := range elem.Attributes {
						if a.Key == "#label" || a.Key == "v-slot:label" {
							var sb strings.Builder
							for _, inner := range elem.Children {
								sb.WriteString(strings.TrimSpace(inner.Dump()))
							}
							labelText = sb.String()
							goto SkipChild
						}
						// Remove hint slot.
						if a.Key == "v-slot:hint" || a.Key == "#hint" {
							goto SkipChild
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
			}
			node.Children = remainingChildren
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
