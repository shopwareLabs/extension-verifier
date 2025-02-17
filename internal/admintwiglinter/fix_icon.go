package admintwiglinter

import (
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type IconFixer struct{}

func init() {
	AddFixer(IconFixer{})
}

func (i IconFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-icon" {
			errors = append(errors, CheckError{
				Message:    "sw-icon is removed, use mt-icon instead with proper size prop.",
				Severity:   "error",
				Identifier: "sw-icon",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (i IconFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (i IconFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-icon" {
			node.Tag = "mt-icon"
			hasSize := false
			var newAttrs []html.Attribute
			for _, attr := range node.Attributes {
				switch strings.ToLower(attr.Key) {
				case "small":
					// Replace "small" with size="16px"
					newAttrs = append(newAttrs, html.Attribute{
						Key:   "size",
						Value: "16px",
					})
					hasSize = true
				case "large":
					// Replace "large" with size="32px"
					newAttrs = append(newAttrs, html.Attribute{
						Key:   "size",
						Value: "32px",
					})
					hasSize = true
				case "size":
					// keep existing size prop
					newAttrs = append(newAttrs, attr)
					hasSize = true
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			// If no size related prop is set, add default size="24px"
			if !hasSize {
				newAttrs = append(newAttrs, html.Attribute{
					Key:   "size",
					Value: "24px",
				})
			}
			node.Attributes = newAttrs
		}
	})
	return nil
}
