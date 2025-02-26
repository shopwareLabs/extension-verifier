package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExternalLinkFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-external-link>Hello World</sw-external-link>`,
			after:       `<mt-external-link>Hello World</mt-external-link>`,
		},
		{
			description: "remove icon attribute",
			before:      `<sw-external-link icon>Hello World</sw-external-link>`,
			after:       `<mt-external-link>Hello World</mt-external-link>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(ExternalLinkFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
