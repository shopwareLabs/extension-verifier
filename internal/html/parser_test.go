package html

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormattingOfHTML(t *testing.T) {
	swBlock := &ElementNode{
		Tag: "sw-button",
		Attributes: []Attribute{
			{
				Key:   "label",
				Value: "Click me",
			},
			{
				Key:   "variant",
				Value: "primary",
			},
		},
	}

	node := &ElementNode{Tag: "template", Attributes: make([]Attribute, 0), Children: NodeList{swBlock}}

	assert.Equal(t, `<template>
    <sw-button
        label="Click me"
        variant="primary"
    ></sw-button>
</template>`, node.Dump())

	simpleButton := &ElementNode{
		Tag: "sw-button",
		Children: NodeList{
			&RawNode{Text: "Click me"},
		},
	}

	assert.Equal(t, `<sw-button>Click me</sw-button>`, simpleButton.Dump())
}

func TestFormatting(t *testing.T) {
	files, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		t.Run(f.Name(), func(t *testing.T) {
			name := f.Name()

			data, err := os.ReadFile(path.Join("testdata", name))
			if err != nil {
				t.Fatal(err)
			}

			stringData := string(data)
			stringParts := strings.SplitN(stringData, "-----", 2)
			stringParts[0] = strings.TrimRight(stringParts[0], "\n")
			stringParts[1] = strings.TrimLeft(stringParts[1], "\n")

			if len(stringParts) != 2 {
				t.Fatalf("file %s does not contain expected delimiter", name)
			}

			parsed, err := NewParser(stringParts[0])
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, stringParts[1], parsed.Dump())
		})
	}

}

func TestChangeElement(t *testing.T) {
	node, err := NewParser(`<sw-select @update:value="onUpdateValue"/>`)
	assert.NoError(t, err)
	TraverseNode(node, func(n *ElementNode) {
		n.Tag = "mt-select"
		for i, attr := range n.Attributes {
			if attr.Key == "@update:value" {
				n.Attributes[i].Key = "@update:modelValue"
			}
		}
	})
	assert.Equal(t, `<mt-select @update:modelValue="onUpdateValue"/>`, node.Dump())
}

func TestBlockParsing(t *testing.T) {
	input := `{% block name %}{% endblock %}`

	node, err := NewParser(input)
	assert.NoError(t, err)

	assert.Equal(t, input, node.Dump())

	block, ok := node[0].(*TwigBlockNode)
	assert.True(t, ok)
	assert.Equal(t, "name", block.Name)
}
