package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateTokenFromHeader(t *testing.T) {
	testCases := []struct {
		name          string
		header        string
		secret        string
		expectedError bool
	}{
		{
			name:          "Valid token",
			header:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			secret:        "secret",
			expectedError: false,
		},
		{
			name:          "Token signed with different secret",
			header:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			secret:        "secret2",
			expectedError: true,
		},
		{
			name:          "Invalid header",
			header:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			secret:        "secret",
			expectedError: true,
		},
		{
			name:          "Invalid header",
			header:        "Bearer a eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			secret:        "secret",
			expectedError: true,
		},
		{
			name:          "Invalid token format",
			header:        "Bearer token",
			secret:        "secret",
			expectedError: true,
		},
		{
			name:          "Expired token",
			header:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyLCJleHAiOjE2MjMzOTIwNDh9.4-vNwoKKYHTQ3uJCyaDZhjsjdGomgWSTDKRhfRd7vEM",
			secret:        "secret",
			expectedError: true,
		},
		{
			name:          "Empty secret",
			header:        "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.XbPfbIHMI6arZ3Y922BhjWgQzWXcXNrz0ogtVhfEd2o",
			secret:        "",
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			err := ValidateTokenFromHeader(tc.header, tc.secret)
			if tc.expectedError {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}
