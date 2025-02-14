package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwitchFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-switch-field>Hello World</sw-switch-field>`,
			after:       `<mt-switch>Hello World</mt-switch>`,
		},
		{
			description: "replace noMarginTop with removeTopMargin",
			before:      `<sw-switch-field noMarginTop />`,
			after:       `<mt-switch removeTopMargin=""/>`,
		},
		{
			description: "remove size, id, ghostValue, padded and partlyChecked",
			before:      `<sw-switch-field size="small" id="example" ghostValue="true" padded="true" partlyChecked="true" />`,
			after:       `<mt-switch/>`,
		},
		{
			description: "replace value with checked",
			before:      `<sw-switch-field value="true" />`,
			after:       `<mt-switch checked="true"/>`,
		},
		{
			description: "convert label slot to label prop",
			before: `<sw-switch-field><template #label>
        Foobar
    </template></sw-switch-field>`,
			after: `<mt-switch label="Foobar"></mt-switch>`,
		},
		{
			description: "remove hint slot and add comment node",
			before: `<sw-switch-field><template #hint>
        Foobar
    </template></sw-switch-field>`,
			after: `<mt-switch></mt-switch>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(SwitchFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
