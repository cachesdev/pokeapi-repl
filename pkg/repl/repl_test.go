package repl

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCleanInput(t *testing.T) {

	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "  hello  world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		actual := CleanInput(c.input)
		require.Len(t, actual, len(c.expected))
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			require.Equal(t, word, expectedWord, "should be equal")
		}
	}
}
