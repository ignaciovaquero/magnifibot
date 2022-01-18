package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateOriginURLsFromArray(t *testing.T) {
	testCases := []struct {
		name          string
		urls          []string
		expectedError bool
	}{
		{
			name:          "Valid array of URLs",
			urls:          []string{"https://validurl.com", "http://validurl.com"},
			expectedError: false,
		},
		{
			name:          "Validate all origins",
			urls:          []string{"*"},
			expectedError: false,
		},
		{
			name:          "Empty array of URLs",
			urls:          []string{},
			expectedError: false,
		},
		{
			name:          "Invalid URL in array",
			urls:          []string{"http://invalid url"},
			expectedError: true,
		},
		{
			name:          "URL with no protocol",
			urls:          []string{"example.com"},
			expectedError: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := ValidateOriginURLsFromArray(tc.urls)
			if tc.expectedError {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
