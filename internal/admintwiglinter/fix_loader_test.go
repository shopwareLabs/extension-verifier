package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoaderFixer(t *testing.T) {
	cases := []struct {
		before string
		after  string
	}{
		{
			before: `<sw-loader />`,
			after:  `<mt-loader/>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(LoaderFixer{}, c.before)
		assert.NoError(t, err)
		assert.Equal(t, c.after, newStr)
	}
}
