package admintwiglinter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCardFixer(t *testing.T) {
	cases := []struct {
		description string
		before      string
		after       string
	}{
		{
			description: "basic component replacement",
			before:      `<sw-card>Hello World</sw-card>`,
			after:       `<mt-card>Hello World</mt-card>`,
		},
		{
			description: "remove contentPadding property",
			before:      `<sw-card contentPadding="true">Hello World</sw-card>`,
			after:       `<mt-card>Hello World</mt-card>`,
		},
		{
			description: "convert aiBadge property to title slot",
			before:      `<sw-card aiBadge>Hello World</sw-card>`,
			after: `<mt-card>
    <slot name="title">
        <sw-ai-copilot-badge></sw-ai-copilot-badge>
    </slot>
    Hello World
</mt-card>`,
		},
	}

	for _, c := range cases {
		newStr, err := runFixerOnString(CardFixer{}, c.before)
		assert.NoError(t, err, c.description)
		assert.Equal(t, c.after, newStr, c.description)
	}
}
