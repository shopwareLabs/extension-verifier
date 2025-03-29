package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
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
				Line:       node.Line,
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
			var newAttrs html.NodeList
			// Process attribute conversions.
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case ":value":
						newAttrs = append(newAttrs, html.Attribute{Key: ":checked", Value: attr.Value})
					case "v-model", "v-model:value":
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
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}
			node.Attributes = newAttrs

			// Process children for slot conversion.
			var labelText string
			var remainingChildren html.NodeList
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					// Handle label slot.
					for _, a := range elem.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							if attr.Key == "#label" || attr.Key == "v-slot:label" {
								var sb strings.Builder
								for _, inner := range elem.Children {
									sb.WriteString(strings.TrimSpace(inner.Dump(0)))
								}
								labelText = sb.String()
								goto SkipChild
							}
							// Remove hint slot.
							if attr.Key == "v-slot:hint" || attr.Key == "#hint" {
								goto SkipChild
							}
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
