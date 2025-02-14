package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextarea(t *testing.T) {
	cases := []struct {
		before string
		after  string
	}{
		{
			before: `<sw-textarea-field></sw-textarea-field>`,
			after:  `<mt-textarea></mt-textarea>`,
		},
		{
			before: `<sw-textarea-field><template #label>FOO</template></sw-textarea-field>`,
			after:  `<mt-textarea label="FOO"></mt-textarea>`,
		},
	}

	for _, c := range cases {
		new, err := runFixerOnString(TextareaFieldFixer{}, c.before)

		assert.NoError(t, err)

		assert.Equal(t, c.after, new)
	}
}
