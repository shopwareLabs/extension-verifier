package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatepickerFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-datepicker :modelValue="myValue" v-model="myValue" @update:modelValue="handler"/>`,
			after: `<mt-datepicker
	:modelValue="myValue"
	v-model="myValue"
	@update:modelValue="handler"
/>`,
		},
		{
			description: "convert label slot to prop",
			before:      `<sw-datepicker><template #label>My Label</template></sw-datepicker>`,
			after:       `<mt-datepicker label="My Label"></mt-datepicker>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(DatepickerFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
