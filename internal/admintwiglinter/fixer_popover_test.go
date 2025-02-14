package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPopover(t *testing.T) {
	cases := []struct {
		before string
		after  string
	}{
		{
			before: `<sw-popover></sw-popover>`,
			after:  `<mt-floating-ui :isOpened="true"></mt-floating-ui>`,
		},
		{
			before: `<sw-popover v-if="bla"></sw-popover>`,
			after:  `<mt-floating-ui :isOpened="bla"></mt-floating-ui>`,
		},
		{
			before: `<sw-popover :zIndex="123"></sw-popover>`,
			after:  `<mt-floating-ui :isOpened="true"></mt-floating-ui>`,
		},
		{
			before: `<sw-popover :resizeWidth="123"></sw-popover>`,
			after:  `<mt-floating-ui :isOpened="true"></mt-floating-ui>`,
		},
	}

	for _, c := range cases {
		new, err := runFixerOnString(PopoverFixer{}, c.before)

		assert.NoError(t, err)

		assert.Equal(t, c.after, new)
	}
}
