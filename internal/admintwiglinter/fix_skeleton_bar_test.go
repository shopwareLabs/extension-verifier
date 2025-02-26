package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSkeletonBarFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			before: `<sw-skeleton-bar>Hello World</sw-skeleton-bar>`,
			after:  `<mt-skeleton-bar>Hello World</mt-skeleton-bar>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(SkeletonBarFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
