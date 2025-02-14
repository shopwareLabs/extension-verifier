package admintwiglinter

import (
	"encoding/json"
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type TabsFixer struct{}

func init() {
	AddFixer(TabsFixer{})
}

func (t TabsFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-tabs" {
			errs = append(errs, CheckError{
				Message:    "sw-tabs is removed, use mt-tabs instead. Review conversion for slots and properties.",
				Severity:   "error",
				Identifier: "sw-tabs",
			})
		}
	})
	return errs
}

func (t TabsFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (t TabsFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-tabs" {
			node.Tag = "mt-tabs"
			var newAttrs []html.Attribute
			// Process attribute conversions.
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "is-vertical":
					newAttrs = append(newAttrs, html.Attribute{Key: "vertical", Value: attr.Value})
				case "align-right":
					// Remove align-right.
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process children for slot conversion.
			var defaultItems []map[string]string
			contentSlotFound := false
			var remainingChildren []html.Node

			for _, child := range node.Children {
				if tpl, ok := child.(*html.ElementNode); ok && tpl.Tag == "template" {
					for _, a := range tpl.Attributes {
						// Process default slot.
						if a.Key == "#default" || a.Key == "v-slot:default" {
							// Find all sw-tabs-item elements in the template.
							for _, itemNode := range tpl.Children {
								if itemElem, ok := itemNode.(*html.ElementNode); ok && itemElem.Tag == "sw-tabs-item" {
									var itemLabel string
									var itemName string
									for _, attr := range itemElem.Attributes {
										if attr.Key == "name" {
											itemName = attr.Value
										}
									}
									// Get inner text for label.
									var sb strings.Builder
									for _, inner := range itemElem.Children {
										sb.WriteString(strings.TrimSpace(inner.Dump()))
									}
									itemLabel = sb.String()
									defaultItems = append(defaultItems, map[string]string{"label": itemLabel, "name": itemName})
								}
							}
							// Skip this template.
							goto NextChild
						}
						// Process content slot.
						if a.Key == "#content" || a.Key == "v-slot:content" {
							contentSlotFound = true
							// Skip content slot.
							goto NextChild
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			NextChild:
			}
			node.Children = remainingChildren

			// If default slot was found, add items prop.
			if len(defaultItems) > 0 {
				if bytes, err := json.Marshal(defaultItems); err == nil {
					// Use single quotes inside the JSON string by replacing double quotes.
					itemsVal := strings.ReplaceAll(string(bytes), "\"", "'")
					node.Attributes = append(node.Attributes, html.Attribute{
						Key:   ":items",
						Value: itemsVal,
					})
				}
			}

			// If a content slot was found, add the event for active item.
			if contentSlotFound {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "@new-item-active",
					Value: "setActiveItem",
				})
			}
		}
	})
	return nil
}
