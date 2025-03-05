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
		{
			description: "elements with block",
			before:      `{% block foo %}<sw-button>Click me</sw-button>{% endblock %}`,
			after: `{% block foo %}
    <sw-button>Click me</sw-button>
{% endblock %}`,
		},
		{
			description: "multi line breaks get removed",
			before: `{% block test %}<sw-button>Click me</sw-button>


<sw-button>Click me</sw-button>{% endblock %}`,
			after: `{% block test %}
    <sw-button>Click me</sw-button>

    <sw-button>Click me</sw-button>
{% endblock %}`,
		},
		{
			description: "multi line between elements only one",
			before:      `<template><foo><bar/></foo></template>`,
			after: `<template>
    <foo>
        <bar/>
    </foo>
</template>`,
		},
		{
			description: "multi line between only elements",
			before: `<template>


<foo>
	<bar/>
</foo>


</template>`,
			after: `<template>
    <foo>
        <bar/>
    </foo>
</template>`,
		},
		{
			description: "long attribute is on new line",
			before:      `<sw-button link="{ name: 'sw.product.detail.pseudovariants', params: { productId: product.id } }"/>`,
			after: `<sw-button
    link="{ name: 'sw.product.detail.pseudovariants', params: { productId: product.id } }"
/>`,
		},
		{
			description: "html element with content gets correct formatting",
			before:      `<template><router-link>{{ item.mainPseudovariant.product.translated.name }}</router-link></template>`,
			after: `<template>
    <router-link>
        {{ item.mainPseudovariant.product.translated.name }}
    </router-link>
</template>`,
		},
		{
			description: "multiple template elements should have a newline between them",
			before:      `<template><div>Template 1</div></template><template><div>Template 2</div></template>`,
			after: `<template>
    <div>Template 1</div>
</template>

<template>
    <div>Template 2</div>
</template>`,
		},
		{
			description: "multiple template elements should have a newline between them with root element",
			before:      `<sw-page><template><div>Template 1</div></template><template><div>Template 2</div></template></sw-page>`,
			after: `<sw-page>
    <template>
        <div>Template 1</div>
    </template>

    <template>
        <div>Template 2</div>
    </template>
</sw-page>`,
		},
		{
			description: "starting tag in html node",
			before:      "<p>{{ $tc('swag-customized-products.detail.tabGeneral.cardExclusion.emptyTitle', (searchTerm.length <= 0) ? 1 : 0) }}</p>",
			after: `<p>
    {{ $tc('swag-customized-products.detail.tabGeneral.cardExclusion.emptyTitle', (searchTerm.length <= 0) ? 1 : 0) }}
</p>`,
		},
		{
			description: "template expression in div",
			before:      "<div>{{ someVariable }}</div>",
			after:       "<div>{{ someVariable }}</div>",
		},
		{
			description: "multiple template expressions",
			before:      "<div>{{ firstVar }}{{ secondVar }}</div>",
			after:       "<div>{{ firstVar }}{{ secondVar }}</div>",
		},
		{
			description: "template expression with text",
			before:      "<div>Before {{ expression }} After</div>",
			after:       "<div>Before {{ expression }} After</div>",
		},
		{
			description: "template expression in nested elements",
			before:      "<div><span>{{ nestedExpression }}</span></div>",
			after: `<div>
    <span>{{ nestedExpression }}</span>
</div>`,
		},
		{
			description: "template expression in router-link with long expression",
			before:      "<router-link>{{ item.mainPseudovariant.product.translated.name }}</router-link>",
			after: `<router-link>
    {{ item.mainPseudovariant.product.translated.name }}
</router-link>`,
		},
		{
			description: "multiple long template expressions",
			before:      "<div>{{ item.mainPseudovariant.product.translated.name }}{{ item.mainPseudovariant.product.translated.description }}</div>",
			after: `<div>
    {{ item.mainPseudovariant.product.translated.name }}
    {{ item.mainPseudovariant.product.translated.description }}
</div>`,
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
			if attr.Key == "@update:value" {
				n.Attributes[i].Key = "@update:modelValue"
			}
		}
	})
	assert.Equal(t, `<mt-select @update:modelValue="onUpdateValue"/>`, node.Dump())
}

func TestMultipleProcessDoesNotChangeFormatting(t *testing.T) {
	code := `{% block sw_import_export_tabs_profiles %}
    {% parent() %}

    <sw-tabs-item :route="{ name: 'iwvs.import.export.index.color' }">
        {{ $tc('iwvs-import-export.page.colorTab') }}
    </sw-tabs-item>
{% endblock %}`

	nodes, err := NewParser(code)
	assert.NoError(t, err)
	assert.Equal(t, code, nodes.Dump())

	nodes, err = NewParser(nodes.Dump())
	assert.NoError(t, err)
	assert.Equal(t, code, nodes.Dump())
}

func TestNoEndingNewlines(t *testing.T) {
	code := `<sw-card-view>
    <sw-card
        class="frosh-tools-tab-scheduled__tasks-card"
        :title="$tc('frosh-tools.tabs.scheduledTaskOverview.title')"
        :isLoading="isLoading"
        :large="true"
        positionIdentifier="frosh-tools-tab-scheduled"
    >
        <template #toolbar>
            <sw-button
                variant="ghost"
                @click="refresh"
            >
                <sw-icon
                    :small="true"
                    name="regular-undo"
                ></sw-icon>
            </sw-button>
            <sw-button
                variant="primary"
                @click="registerScheduledTasks"
            >
                {{ $tc('frosh-tools.scheduledTasksRegisterStarted') }}
            </sw-button>
        </template>
        <sw-entity-listing
            :showSelection="false"
            :fullPage="false"
            :allowInlineEdit="true"
            :allowEdit="false"
            :allowDelete="false"
            :showActions="true"
            :repository="scheduledRepository"
            :items="items"
            :columns="columns"
        >
            <template #column-lastExecutionTime="{ item }">
                {{ date(item.lastExecutionTime, {hour: '2-digit', minute: '2-digit'}) }}
            </template>

            <template
                #column-nextExecutionTime="{ item, column, compact, isInlineEdit }"
            >
                <sw-data-grid-inline-edit
                    v-if="isInlineEdit"
                    :column="column"
                    :compact="compact"
                    :value="item[column.property]"
                    @update:value="item[column.property] = $event"
                >
                </sw-data-grid-inline-edit>
                <span v-else>
                    {{ date(item.nextExecutionTime, {hour: '2-digit', minute: '2-digit'}) }}
                </span>
            </template>

            <template #actions="{ item }">
                <sw-context-menu-item
                    variant="primary"
                    @click="runTask(item)"
                >
                    {{ $tc('frosh-tools.runManually') }}
                </sw-context-menu-item>
                <sw-context-menu-item
                    variant="primary"
                    @click="scheduleTask(item)"
                >
                    {{ $tc('frosh-tools.setToScheduled') }}
                </sw-context-menu-item>
                <sw-context-menu-item
                    variant="primary"
                    @click="scheduleTask(item, true)"
                >
                    {{ $tc('frosh-tools.setToScheduledImmediately') }}
                </sw-context-menu-item>
            </template>
        </sw-entity-listing>
    </sw-card>
    <sw-modal
        v-if="taskError"
        :title="$tc('global.default.error')"
        @modal-close="taskError = null"
    >
        <pre
            v-if="typeof taskError === 'object'"
            v-text="taskError"
        />
        <div
            v-else
            v-html="taskError"
        />
        <template #modal-footer>
            <sw-button
                size="small"
                @click="taskError = null"
            >
                {{ $tc('global.default.close') }}
            </sw-button>
        </template>
    </sw-modal>
</sw-card-view>
`

	nodes, err := NewParser(code)
	assert.NoError(t, err)
	assert.Equal(t, code, nodes.Dump())
}
