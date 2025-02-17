package admintwiglinter

import (
	"fmt"
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type ButtonFixer struct{}

func init() {
	AddFixer(ButtonFixer{})
}

func (b ButtonFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-button" {
			errors = append(errors, CheckError{
				Message:    "sw-button is removed, use mt-button instead. Please review conversion for variant and router-link.",
				Severity:   "error",
				Identifier: "sw-button",
				Line:       node.Line,
			})
		}
	})
	return errors
}

func (b ButtonFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (b ButtonFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-button" {
			node.Tag = "mt-button"
			var newAttrs []html.Attribute
			// Flags to determine additional properties.
			addGhost := false
			for _, attr := range node.Attributes {
				switch attr.Key {
				case "variant":
					lower := strings.ToLower(attr.Value)
					switch lower {
					case "ghost":
						// Remove variant and set ghost.
						addGhost = true
					case "danger":
						// Change value to critical.
						attr.Value = "critical"
						newAttrs = append(newAttrs, attr)
					case "ghost-danger":
						// Set critical and also ghost.
						attr.Value = "critical"
						newAttrs = append(newAttrs, attr)
						addGhost = true
					case "contrast", "context":
						// Remove attribute
					default:
						newAttrs = append(newAttrs, attr)
					}
				case "router-link":
					// Replace with @click event.
					val := attr.Value
					newAttrs = append(newAttrs, html.Attribute{
						Key:   "@click",
						Value: fmt.Sprintf("this.$router.push('%s')", val),
					})
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			if addGhost {
				newAttrs = append(newAttrs, html.Attribute{
					Key: "ghost",
				})
			}
			node.Attributes = newAttrs
		}
	})
	return nil
}
