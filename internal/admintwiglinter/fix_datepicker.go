package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
)

type DatepickerFixer struct{}

func init() {
	AddFixer(DatepickerFixer{})
}

func (d DatepickerFixer) Check(nodes []html.Node) []CheckError {
	var checkErrors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-datepicker" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-datepicker is removed, use mt-datepicker instead. Please review the conversion for the label property.",
				Severity:   "error",
				Identifier: "sw-datepicker",
				Line:       node.Line,
			})
		}
	})
	return checkErrors
}

func (d DatepickerFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (d DatepickerFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-datepicker" {
			node.Tag = "mt-datepicker"

			var newAttrs html.NodeList
			// Update attribute names.
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case ":value":
						attr.Key = ":model-value"
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
				} else {
					// If it's not an html.Attribute (e.g., TwigIfNode), preserve it as is
					newAttrs = append(newAttrs, attrNode)
				}
			}
			node.Attributes = newAttrs

			// Convert label slot to label property.
			label := ""
			var remainingChildren html.NodeList
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					if elem.Tag == "template" {
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
