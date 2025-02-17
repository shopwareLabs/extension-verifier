package admintwiglinter

import (
	"github.com/shopware/extension-verifier/internal/html"
	"github.com/shopware/shopware-cli/version"
)

type LoaderFixer struct{}

func init() {
	AddFixer(LoaderFixer{})
}

func (l LoaderFixer) Check(nodes []html.Node) []CheckError {
	var errs []CheckError
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-loader" {
			errs = append(errs, CheckError{
				Message:    "sw-loader is removed, use mt-loader instead.",
				Severity:   "error",
				Identifier: "sw-loader",
				Line:       node.Line,
			})
		}
	})
	return errs
}

func (l LoaderFixer) Supports(v *version.Version) bool {
	return shopware67Constraint.Check(v)
}

func (l LoaderFixer) Fix(nodes []html.Node) error {
	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-loader" {
			node.Tag = "mt-loader"
		}
	})
	return nil
}
