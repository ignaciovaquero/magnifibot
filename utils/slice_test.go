package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAllButOne(t *testing.T) {
	testCases := []struct {
		name     string
		items    []string
		item     string
		expected []string
	}{
		{
			name:     "Multiple items in array, matching element",
			items:    []string{"one", "two", "three"},
			item:     "two",
			expected: []string{"one", "three"},
		},
		{
			name:     "Multiple items in array, non-matching element",
			items:    []string{"one", "two", "three"},
			item:     "four",
			expected: []string{"one", "two", "three"},
		},
		{
			name:     "Single item in array, matching element",
			items:    []string{"one"},
			item:     "one",
			expected: []string{},
		},
		{
			name:     "No items in array",
			items:    []string{},
			item:     "four",
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			actual := AllButOne(tc.items, tc.item)
			assert.Equal(tt, tc.expected, actual)
		})
	}
}
