package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProgressBarFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-progress-bar/>`,
			after:       `<mt-progress-bar/>`,
		},
		{
			description: "replace value with modelValue",
			before:      `<sw-progress-bar value="5"/>`,
			after:       `<mt-progress-bar modelValue="5"/>`,
		},
		{
			description: "replace v-model:value with v-model",
			before:      `<sw-progress-bar v-model:value="myValue"/>`,
			after:       `<mt-progress-bar v-model="myValue"/>`,
		},
		{
			description: "replace update:value event",
			before:      `<sw-progress-bar @update:value="updateValue"/>`,
			after:       `<mt-progress-bar @update:modelValue="updateValue"/>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(ProgressBarFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
