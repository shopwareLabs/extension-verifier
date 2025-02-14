package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectFieldFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-select-field />`,
			after:       `<mt-select/>`,
		},
		{
			description: "replace value with modelValue",
			before:      `<sw-select-field :value="selectedValue" />`,
			after:       `<mt-select :modelValue="selectedValue"/>`,
		},
		{
			description: "replace v-model:value with v-model",
			before:      `<sw-select-field v-model:value="selectedValue" />`,
			after:       `<mt-select v-model="selectedValue"/>`,
		},
		{
			description: "convert options prop format",
			before:      `<sw-select-field :options="[ { name: 'Option 1', id: 1 }, { name: 'Option 2', id: 2 } ]" />`,
			after:       `<mt-select :options="[ { label: 'Option 1', value: 1 }, { label: 'Option 2', value: 2 } ]"/>`,
		},
		{
			description: "remove aside prop",
			before:      `<sw-select-field :aside="true" />`,
			after:       `<mt-select/>`,
		},
		{
			description: "convert default slot with option children to options prop",
			before: `<sw-select-field>
    <option value="1">Option 1</option>
    <option value="2">Option 2</option>
</sw-select-field>`,
			after: `<mt-select :options="[{"label":"Option 1","value":"1"},{"label":"Option 2","value":"2"}]"></mt-select>`,
		},
		{
			description: "convert label slot to label prop",
			before: `<sw-select-field>
    <template #label>
        My Label
    </template>
</sw-select-field>`,
			after: `<mt-select label="My Label"></mt-select>`,
		},
		{
			description: "replace update:value event with update:modelValue",
			before:      `<sw-select-field @update:value="onUpdateValue" />`,
			after:       `<mt-select @update:modelValue="onUpdateValue"/>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(SelectFieldFixer{}, c.before)
		assert.NoError(t, err, c.description)
		// Normalize whitespace for comparison.
		assert.Equal(t, c.after, newStr, c.description)
	}
}
