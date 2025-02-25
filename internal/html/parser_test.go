package html

import (
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

func TestParseAndPrint(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic element",
			before:      `<sw-button>Click me</sw-button>`,
			after:       `<sw-button>Click me</sw-button>`,
		},
		{
			description: "sub-nodes",
			before:      `<template><div><sw-button>Foo</sw-button></div></template>`,
			after: `<template>
	<div>
		<sw-button>Foo</sw-button>
	</div>
</template>`,
		},
		{
			description: "attributes single",
			before:      `<sw-button variant="primary">Click me</sw-button>`,
			after:       `<sw-button variant="primary">Click me</sw-button>`,
		},
		{
			description: "attributes",
			before:      `<sw-button variant="primary" foo="bla">Click me</sw-button>`,
			after: `<sw-button
	variant="primary"
	foo="bla"
>Click me</sw-button>`,
		},
		{
			description: "children with comment",
			before:      `<sw-button><!-- comment --></sw-button>`,
			after:       `<sw-button><!-- comment --></sw-button>`,
		},
		{
			description: "multiple comments",
			before:      `<div><!-- header -->Content<!-- footer --></div>`,
			after:       `<div><!-- header -->Content<!-- footer --></div>`,
		},
		{
			description: "comment with nested tags",
			before:      `<!-- <div>this is commented out</div> --><div>actual content</div>`,
			after:       `<!-- <div>this is commented out</div> --><div>actual content</div>`,
		},
		{
			description: "comment with special characters",
			before:      `<div><!-- special chars: & < > " ' --></div>`,
			after:       `<div><!-- special chars: & < > " ' --></div>`,
		},
	}

	for _, c := range cases {
		node, err := NewParser(c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, node.Dump(), c.description)
	}
}

func TestChangeElement(t *testing.T) {
	node, err := NewParser(`<sw-select @update:value="onUpdateValue"/>`)

	assert.NoError(t, err)

	TraverseNode(node, func(n *ElementNode) {
		n.Tag = "mt-select"
		for i, attr := range n.Attributes {
			if attr.Key == "update:value" {
				n.Attributes[i].Key = "update:modelValue"
			}
		}
	})

	assert.Equal(t, `<mt-select @update:modelValue="onUpdateValue"/>`, node.Dump())
}
