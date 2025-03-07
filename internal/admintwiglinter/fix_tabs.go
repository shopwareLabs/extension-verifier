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
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-tabs" {
			errs = append(errs, CheckError{
				Message:    "sw-tabs is removed, use mt-tabs instead. Review conversion for slots and properties.",
				Severity:   "error",
				Identifier: "sw-tabs",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (t TabsFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (t TabsFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-tabs" {
			node.Tag = "mt-tabs"
			var newAttrs html.NodeList
			// Process attribute conversions.
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case "is-vertical":
						newAttrs = append(newAttrs, html.Attribute{Key: "vertical", Value: attr.Value})
					case "align-right":
						// Remove align-right.
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
			var defaultItems []map[string]string
			contentSlotFound := false
			var remainingChildren html.NodeList

			for _, child := range node.Children {
				if tpl, ok := child.(*html.ElementNode); ok && tpl.Tag == "template" {
					for _, a := range tpl.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							// Process default slot.
							if attr.Key == "#default" || attr.Key == "v-slot:default" {
								// Find all sw-tabs-item elements in the template.
								for _, itemNode := range tpl.Children {
									if itemElem, ok := itemNode.(*html.ElementNode); ok && itemElem.Tag == "sw-tabs-item" {
										var itemLabel string
										var itemName string
										for _, itemAttr := range itemElem.Attributes {
											if attr, ok := itemAttr.(html.Attribute); ok {
												if attr.Key == "name" {
													itemName = attr.Value
												}
											}
										}
										// Get inner text for label.
										var sb strings.Builder
										for _, inner := range itemElem.Children {
											sb.WriteString(strings.TrimSpace(inner.Dump(0)))
										}
										itemLabel = sb.String()
										defaultItems = append(defaultItems, map[string]string{"label": itemLabel, "name": itemName})
									}
								}
								// Skip this template.
								goto SkipChild
							}
							// Process content slot.
							if attr.Key == "#content" || attr.Key == "v-slot:content" {
								contentSlotFound = true
								// Skip content slot.
								goto SkipChild
							}
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
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
