package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
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

			for i, attr := range node.Attributes {
				switch attr.Key {
				case "value":
					node.Attributes[i].Key = "modelValue"
				case "v-model:value":
					node.Attributes[i].Key = "v-model"
				case "update:value":
					node.Attributes[i].Key = "update:modelValue"
				}
			}

			label := ""

			for _, children := range node.Children {
				if element, ok := children.(*html.ElementNode); ok {
					if element.Tag == "template" {
						for _, attr := range element.Attributes {
							if attr.Key == "#label" {
								for _, child := range element.Children {
									label = label + child.Dump()
								}
							}
						}
					}
				}
			}

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
