package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTabsFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "convert default slot to items prop",
			before: `<sw-tabs><template #default>
        <sw-tabs-item name="tab1">Tab 1</sw-tabs-item>
        <sw-tabs-item name="tab2">Tab 2</sw-tabs-item>
    </template></sw-tabs>`,
			after: `<mt-tabs :items="[{'label':'Tab 1','name':'tab1'},{'label':'Tab 2','name':'tab2'}]"></mt-tabs>`,
		},
		{
			description: "remove content slot and add event",
			before: `<sw-tabs><template #content="{ active }">
        The current active item is {{ active }}
    </template></sw-tabs>`,
			after: `<mt-tabs @new-item-active="setActiveItem"></mt-tabs>`,
		},
		{
			description: "rename is-vertical to vertical and remove align-right",
			before:      `<sw-tabs is-vertical align-right />`,
			after:       `<mt-tabs vertical/>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(TabsFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
