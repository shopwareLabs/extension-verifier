package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDatepickerFixer(t *testing.T) {
	cases := []struct {
		before string
		after  string
	}{
		{
			before: `<sw-datepicker></sw-datepicker>`,
			after:  `<mt-datepicker></mt-datepicker>`,
		},
		{
			before: `<sw-datepicker :value="myValue" v-model:value="myValue" @update:value="handler"></sw-datepicker>`,
			after:  `<mt-datepicker :modelValue="myValue" v-model="myValue" @update:modelValue="handler"></mt-datepicker>`,
		},
		{
			before: `<sw-datepicker><template #label>My Label</template></sw-datepicker>`,
			after:  `<mt-datepicker label="My Label"></mt-datepicker>`,
		},
	}

	for _, c := range cases {
		new, err := runFixerOnString(DatepickerFixer{}, c.before)
		assert.NoError(t, err)
		assert.Equal(t, c.after, new)
	}
}
