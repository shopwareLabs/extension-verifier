package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type PasswordFieldFixer struct{}

func init() {
	AddFixer(PasswordFieldFixer{})
}

func (p PasswordFieldFixer) Check(nodes []html.Node) []CheckError {
	var checkErrors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-password-field" {
			checkErrors = append(checkErrors, CheckError{
				Message:    "sw-password-field is removed, use mt-password-field instead. Please review conversion for label/hint properties.",
				Severity:   "error",
				Identifier: "sw-password-field",
				Line:       node.Line,
			})
		}
	})
	return checkErrors
}

func (p PasswordFieldFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (p PasswordFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-password-field" {
			node.Tag = "mt-password-field"

			// Update or remove attributes
			var newAttrs []html.Attribute
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "value":
					attr.Key = "modelValue"
					newAttrs = append(newAttrs, attr)
				case "v-model:value":
					attr.Key = "v-model"
					newAttrs = append(newAttrs, attr)
				case "size":
					if attr.Value == "medium" {
						attr.Value = "default"
					}
					newAttrs = append(newAttrs, attr)
				case "isInvalid":
					// remove attribute
				case "@update:value":
					attr.Key = "@update:modelValue"
					newAttrs = append(newAttrs, attr)
				case "@base-field-mounted":
					// remove attribute
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process slot children for label and hint
			var label, hint string
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					for _, attr := range elem.Attributes {
						if attr.Key == "#label" {
							var content string
							for _, inner := range elem.Children {
								content += strings.TrimSpace(inner.Dump())
							}
							label = strings.Replace(content, "Label", "label", 1)
						}
						if attr.Key == "#hint" {
							var content string
							for _, inner := range elem.Children {
								content += strings.TrimSpace(inner.Dump())
							}
							hint = strings.Replace(content, "Hint", "hint", 1)
						}
					}
				}
			}
			// Remove original children after processing slots
			node.Children = []html.Node{}
			if label != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: label,
				})
			}
			if hint != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "hint",
					Value: hint,
				})
			}
		}
	})
	return nil
}
