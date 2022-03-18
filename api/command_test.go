package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		name     string
		command  Command
		expected bool
	}{
		{
			name:     "suscribe command",
			command:  ToCommand("suscribirme"),
			expected: true,
		},
		{
			name:     "unsuscribe command",
			command:  ToCommand("baja"),
			expected: true,
		},
		{
			name:     "invalid suscribe command",
			command:  ToCommand("suscribe"),
			expected: false,
		},
		{
			name:     "invalid unsuscribe command",
			command:  ToCommand("unsuscribe"),
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			actual := test.command.IsValid()
			assert.Equal(tt, test.expected, actual)
		})
	}
}

func TestGetValidCommands(t *testing.T) {
	tests := []struct {
		name     string
		expected []Command
	}{
		{
			name:     "Get valid commands",
			expected: []Command{"/suscribirme", "/baja"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			actual := GetValidCommands()
			assert.ElementsMatch(tt, test.expected, actual)
		})
	}
}

func TestGetValidCommandsString(t *testing.T) {
	tests := []struct {
		name     string
		expected []string
	}{
		{
			name:     "",
			expected: []string{"/suscribirme", "/baja"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			actual := GetValidCommandsString()
			assert.ElementsMatch(tt, test.expected, actual)
		})
	}
}
