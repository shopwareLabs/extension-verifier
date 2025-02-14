package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type SkeletonBarFixer struct{}

func init() {
	AddFixer(SkeletonBarFixer{})
}

func (s SkeletonBarFixer) Check(nodes []html.Node) []CheckError {
	var errors []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-skeleton-bar" {
			errors = append(errors, CheckError{
				Message:    "sw-skeleton-bar is removed, use mt-skeleton-bar instead.",
				Severity:   "error",
				Identifier: "sw-skeleton-bar",
			})
		}
	})
	return errors
}

func (s SkeletonBarFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (s SkeletonBarFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-skeleton-bar" {
			node.Tag = "mt-skeleton-bar"
		}
	})
	return nil
}
