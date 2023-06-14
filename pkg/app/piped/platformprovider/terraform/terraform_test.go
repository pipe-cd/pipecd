package terraform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPlanHasChangeRegex(t *testing.T) {
	t.Parallel()

	testcases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "older than v1.5.0",
			input:    "Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 0 to import, 1 to add, 2 to change, 3 to destroy.", "1", "2", "3"},
		},
		{
			name:     "later than v1.5.0",
			input:    "Plan: 1 to add, 2 to change, 3 to destroy.",
			expected: []string{"Plan: 1 to add, 2 to change, 3 to destroy.", "1", "2", "3"},
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.expected, planHasChangeRegex.FindStringSubmatch(tc.input))
		})
	}
}
