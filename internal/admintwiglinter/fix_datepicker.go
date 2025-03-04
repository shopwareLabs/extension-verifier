package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
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

			// Update attribute names.
			for i, attr := range node.Attributes {
				switch attr.Key {
				case ":value":
					node.Attributes[i].Key = ":model-value"
				case "v-model:value":
					node.Attributes[i].Key = "v-model"
				case "@update:value":
					node.Attributes[i].Key = "@update:model-value"
				}
			}

			label := ""
			// Convert label slot to label property.
			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					if elem.Tag == "template" {
						for _, attr := range elem.Attributes {
							if attr.Key == "#label" {
								for _, inner := range elem.Children {
									label += inner.Dump()
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
