package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type CardFixer struct{}

func init() {
	AddFixer(CardFixer{})
}

func (c CardFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-card" {
			errors = append(errors, CheckError{
				Message:    "sw-card is removed, use mt-card instead. Review conversion for aiBadge and contentPadding.",
				Severity:   "error",
				Identifier: "sw-card",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (c CardFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (c CardFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-card" {
			node.Tag = "mt-card"
			var newAttrs []html.Attribute
			aiBadgeFound := false
			// Process attributes: remove aiBadge and contentPadding.
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "aiBadge", "contentPadding":
					if attr.Key == "aiBadge" {
						aiBadgeFound = true
					}
					// remove attribute
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// If aiBadge was present, add title slot with sw-ai-copilot-badge.
			if aiBadgeFound {
				aiBadgeSlot := &html.ElementNode{
					Tag: "slot",
					Attributes: []html.Attribute{
						{Key: "name", Value: "title"},
					},
					Children: []html.Node{
						&html.ElementNode{Tag: "sw-ai-copilot-badge"},
					},
				}
				// Prepend the title slot to existing children.
				node.Children = append([]html.Node{aiBadgeSlot}, node.Children...)
			}
		}
	})
	return nil
}
