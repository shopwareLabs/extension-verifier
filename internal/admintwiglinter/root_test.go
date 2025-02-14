package admintwiglinter

import (
	"strings"
	"testing"

	"github.com/shopware/extension-verifier/internal/html"
	"github.com/stretchr/testify/assert"
)

func TestXxx(t *testing.T) {
	test := `{% block sw_order_line_items_grid_actions %}
    {% parent %}

    <sw-popover/>

    <zeobv-bundle-connections-modal
        v-if="bundleConnectionModalOpen"
        :item="bundleConnectionModalItem"
        :order="order"
        :isLoading="isLoading"
        :onClose="closeBundleConnectionModal"
    ></zeobv-bundle-connections-modal>
{% endblock %}`

	expected := `{% block sw_order_line_items_grid_actions %}
    {% parent %}

    <mt-floating-ui/>

    <zeobv-bundle-connections-modal v-if="bundleConnectionModalOpen" :item="bundleConnectionModalItem" :order="order" :isLoading="isLoading" :onClose="closeBundleConnectionModal"></zeobv-bundle-connections-modal>
{% endblock %}`

	nodes, err := html.NewParser(test)

	html.TraverseNode(nodes, func(node *html.ElementNode) {
		if node.Tag == "sw-popover" {
			node.Tag = "mt-floating-ui"
		}
	})

	assert.NoError(t, err)

	var buf strings.Builder

	for _, node := range nodes {
		buf.WriteString(node.Dump())
	}

	assert.Equal(t, expected, buf.String())
}
