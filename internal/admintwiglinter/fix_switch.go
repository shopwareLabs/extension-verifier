package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shyim/go-version"
)

type SwitchFixer struct{}

func init() {
	AddFixer(SwitchFixer{})
}

func (s SwitchFixer) Check(nodes []html.Node) []CheckError {
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-switch-field" {
			errs = append(errs, CheckError{
				Message:    "sw-switch-field is removed, use mt-switch instead. Review conversion for props, events and slots.",
				Severity:   "error",
				Identifier: "sw-switch-field",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (s SwitchFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (s SwitchFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-switch-field" {
			node.Tag = "mt-switch"
			var newAttrs html.NodeList
			// Process attribute conversions.
			for _, attrNode := range node.Attributes {
				// Check if the attribute is an html.Attribute
				if attr, ok := attrNode.(html.Attribute); ok {
					switch attr.Key {
					case "noMarginTop":
						newAttrs = append(newAttrs, html.Attribute{Key: "removeTopMargin"})
					case "size", "id", "ghostValue", "padded", "partlyChecked":
						// remove these attributes
					case "value":
						newAttrs = append(newAttrs, html.Attribute{Key: "checked", Value: attr.Value})
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
			var labelText string
			var remainingChildren html.NodeList
			for _, child := range node.Children {
				// Check if child is a slot element.
				if elem, ok := child.(*html.ElementNode); ok && elem.Tag == "template" {
					for _, a := range elem.Attributes {
						if attr, ok := a.(html.Attribute); ok {
							if attr.Key == "#label" {
								var sb strings.Builder
								for _, inner := range elem.Children {
									sb.WriteString(strings.TrimSpace(inner.Dump(0)))
								}
								labelText = sb.String()
								goto SkipChild
							}
							if attr.Key == "#hint" {
								goto SkipChild
							}
						}
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
			}
			// Remove all slot children.
			node.Children = remainingChildren
			// If label slot found, add label attribute.
			if labelText != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: labelText,
				})
			}
		}
	})
	return nil
}
