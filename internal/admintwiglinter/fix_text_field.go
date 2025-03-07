package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type TextFieldFixer struct{}

func init() {
	AddFixer(TextFieldFixer{})
}

func (t TextFieldFixer) Check(nodes []html.Node) []CheckError {
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-text-field" {
			errs = append(errs, CheckError{
				Message:    "sw-text-field is removed, use mt-text-field instead. Review conversion for props, events and label slot.",
				Severity:   "error",
				Identifier: "sw-text-field",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (t TextFieldFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (t TextFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-text-field" {
			node.Tag = "mt-text-field"
			var newAttrs html.NodeList
			// Process attributes conversion.
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case "value":
						attr.Key = "model-value"
						newAttrs = append(newAttrs, attr)
					case "v-model:value":
						attr.Key = "v-model"
						newAttrs = append(newAttrs, attr)
					case "size":
						if attr.Value == "medium" {
							attr.Value = "default"
						}
						newAttrs = append(newAttrs, attr)
					case "isInvalid", "aiBadge", "@base-field-mounted":
						// remove these attributes
					case "@update:value":
						attr.Key = "@update:model-value"
						newAttrs = append(newAttrs, attr)
					default:
						newAttrs = append(newAttrs, attr)
					}
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}
			node.Attributes = newAttrs

			// Process label slot: convert <template #label>...</template> to label prop.
			label := ""
			var remainingChildren html.NodeList
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					for _, a := range elem.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							if attr.Key == "#label" {
								var sb strings.Builder
								for _, inner := range elem.Children {
									sb.WriteString(strings.TrimSpace(inner.Dump(0)))
								}
								label = sb.String()
								goto SkipChild
							}
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
			}
			node.Children = remainingChildren
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
