package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
)

type TextareaFieldFixer struct{}

func init() {
	AddFixer(TextareaFieldFixer{})
}

func (t TextareaFieldFixer) Check(nodes []html.Node) []CheckError {
	var checkErrors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-textarea-field" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-textarea-field is removed, use mt-textarea instead. Please manually review the new API differences.",
				Severity:   "error",
				Identifier: "sw-textarea-field",
				Line:       node.Line,
			})
		}
	})
	return checkErrors
}

func (t TextareaFieldFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (t TextareaFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-textarea-field" {
			node.Tag = "mt-textarea"
			var newAttrs html.NodeList

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
					case "update:value":
						attr.Key = "update:model-value"
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

			label := ""
			var remainingChildren html.NodeList

			for _, child := range node.Children {
				if element, ok := child.(*html.ElementNode); ok && element.Tag == "template" {
					for _, a := range element.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							if attr.Key == "#label" {
								var sb strings.Builder
								for _, inner := range element.Children {
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
