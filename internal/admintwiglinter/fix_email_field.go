package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type EmailFieldFixer struct{}

func init() {
	AddFixer(EmailFieldFixer{})
}

func (e EmailFieldFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-email-field" {
			errors = append(errors, CheckError{
				Message:    "sw-email-field is removed, use mt-email-field instead. Review conversion for props, events and label slot.",
				Severity:   "error",
				Identifier: "sw-email-field",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (e EmailFieldFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (e EmailFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-email-field" {
			node.Tag = "mt-email-field"
			var newAttrs []html.Attribute
			for _, attr := range node.Attributes {
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
					// remove attribute
				case "@update:value":
					attr.Key = "@update:model-value"
					newAttrs = append(newAttrs, attr)
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process label slot.
			label := ""
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					for _, a := range elem.Attributes {
						if a.Key == "#label" {
							var content string
							for _, inner := range elem.Children {
								content += strings.TrimSpace(inner.Dump())
							}
							label = content
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
