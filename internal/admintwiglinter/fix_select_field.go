package admintwiglinter

import (
	"encoding/json"
	"strings"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type SelectFieldFixer struct{}

func init() {
	AddFixer(SelectFieldFixer{})
}

func (s SelectFieldFixer) Check(nodes []html.Node) []CheckError {
	// ...existing code...
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-select-field" {
			errs = append(errs, CheckError{
				Message:    "sw-select-field is removed, use mt-select instead. Review conversion for props, slots and events.",
				Severity:   "error",
				Identifier: "sw-select-field",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (s SelectFieldFixer) Supports(v *version.Version) bool {
	// ...existing code...
	return shopware67Constraint.Check(v)
}

func (s SelectFieldFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-select-field" {
			node.Tag = "mt-select"

			var newAttrs []html.Attribute
			// Flag to check if options prop is already set.
			optionsSet := false

			for _, attr := range node.Attributes {
				switch attr.Key {
				case ":value":
					newAttrs = append(newAttrs, html.Attribute{Key: ":model-value", Value: attr.Value})
				case "v-model:value":
					newAttrs = append(newAttrs, html.Attribute{Key: "v-model", Value: attr.Value})
				case ":aside":
					// Remove aside prop.
				case ":options":
					// Convert options format: replace "name" with "label" and "id" with "value"
					converted := strings.ReplaceAll(attr.Value, "name", "label")
					converted = strings.ReplaceAll(converted, "id", "value")
					newAttrs = append(newAttrs, html.Attribute{Key: ":options", Value: converted})
					optionsSet = true
				case "@update:value":
					newAttrs = append(newAttrs, html.Attribute{Key: "@update:model-value", Value: attr.Value})
				default:
					newAttrs = append(newAttrs, attr)
				}
			}
			node.Attributes = newAttrs

			// Process children for slot conversion.
			var labelText string
			var optionObjects []map[string]interface{}
			var remainingChildren []html.Node

			for _, child := range node.Children {
				if elem, ok := child.(*html.ElementNode); ok {
					// Convert label slot to label prop.
					if elem.Tag == "template" {
						for _, a := range elem.Attributes {
							if a.Key == "#label" || a.Key == "v-slot:label" {
								var sb strings.Builder
								for _, inner := range elem.Children {
									sb.WriteString(strings.TrimSpace(inner.Dump()))
								}
								labelText = sb.String()
								goto SkipChild
							}
						}
					}
					// Collect <option> children from default slot.
					if elem.Tag == "option" {
						opt := make(map[string]interface{})
						// Get option value from attributes.
						for _, a := range elem.Attributes {
							if a.Key == "value" {
								opt["value"] = a.Value
							}
						}
						// Get option label from inner text.
						var sb strings.Builder
						for _, inner := range elem.Children {
							sb.WriteString(strings.TrimSpace(inner.Dump()))
						}
						opt["label"] = sb.String()
						optionObjects = append(optionObjects, opt)
						goto SkipChild
					}
				}
				remainingChildren = append(remainingChildren, child)
			SkipChild:
			}
			// Remove all children slots.
			node.Children = remainingChildren

			// If label slot was set, add label attribute.
			if labelText != "" {
				node.Attributes = append(node.Attributes, html.Attribute{
					Key:   "label",
					Value: labelText,
				})
			}

			// If default <option> elements were found and options prop not already set, build options prop.
			if !optionsSet && len(optionObjects) > 0 {
				// Serialize optionObjects slice to JSON-like string.
				bytes, err := json.Marshal(optionObjects)
				if err == nil {
					node.Attributes = append(node.Attributes, html.Attribute{
						Key:   ":options",
						Value: string(bytes),
					})
				}
			}

			node.Children = []html.Node{}
		}
	})
	return nil
}
