package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUrlFieldFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-url-field />`,
			after:       `<mt-url-field/>`,
		},
		{
			description: "replace value with modelValue",
			before:      `<sw-url-field value="Hello World" />`,
			after:       `<mt-url-field modelValue="Hello World"/>`,
		},
		{
			description: "replace v-model:value with v-model",
			before:      `<sw-url-field v-model:value="myValue"/>`,
			after:       `<mt-url-field v-model="myValue"/>`,
		},
		{
			description: "replace update:value event",
			before:      `<sw-url-field @update:value="updateValue" />`,
			after:       `<mt-url-field @update:modelValue="updateValue"/>`,
		},
		{
			description: "process label slot",
			before: `<sw-url-field>
    <template #label>
        My Label
    </template>
</sw-url-field>`,
			after: `<mt-url-field label="My Label"></mt-url-field>`,
		},
		{
			description: "remove hint slot",
			before: `<sw-url-field>
    <template #hint>
        My Hint
    </template>
</sw-url-field>`,
			after: `<mt-url-field></mt-url-field>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(UrlFieldFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
