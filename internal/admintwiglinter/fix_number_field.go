package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type NumberFieldFixer struct{}

func init() {
	AddFixer(NumberFieldFixer{})
}

func (n NumberFieldFixer) Check(nodes []html.Node) []CheckError {
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-number-field" {
			errs = append(errs, CheckError{
				Message:    "sw-number-field is removed, use mt-number-field instead. Please review conversion for props, events and label slot.",
				Severity:   "error",
				Identifier: "sw-number-field",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (n NumberFieldFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (n NumberFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-number-field" {
			node.Tag = "mt-number-field"
			var newAttrs []html.Attribute
			for _, attr := range node.Attributes {
				switch attr.Key {
				case ":value":
					newAttrs = append(newAttrs, html.Attribute{
						Key:   ":model-value",
						Value: attr.Value,
					})
				case "v-model:value":
					newAttrs = append(newAttrs, html.Attribute{
						Key:   ":model-value",
						Value: attr.Value,
					})
					newAttrs = append(newAttrs, html.Attribute{
						Key:   "@change",
						Value: attr.Value + " = $event",
					})
				case "@update:value":
					newAttrs = append(newAttrs, html.Attribute{
						Key:   "@change",
						Value: attr.Value,
					})
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			var label string
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					for _, a := range elem.Attributes {
						if a.Key == "#label" {
							var sb strings.Builder
							for _, inner := range elem.Children {
								sb.WriteString(strings.TrimSpace(inner.Dump()))
							}
							label = sb.String()
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
