package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/mymmrac/telego"
	"github.com/stretchr/testify/assert"
)

type MockTelegram struct {
	messageID int
	err       error
}

func (m *MockTelegram) SendMessage(params *telego.SendMessageParams) (*telego.Message, error) {
	return &telego.Message{
		MessageID: m.messageID,
	}, m.err
}

func TestSetTelegram(t *testing.T) {
	tests := []struct {
		name,
		chatID,
		message string
		messageID,
		expected int
		err           error
		errorExpected bool
	}{
		{
			name:          "valid message",
			chatID:        "12",
			message:       "message",
			messageID:     12,
			expected:      12,
			err:           nil,
			errorExpected: false,
		},
		{
			name:          "invalid chat ID",
			chatID:        "id",
			message:       "message",
			messageID:     1,
			expected:      0,
			err:           nil,
			errorExpected: true,
		},
		{
			name:          "empty chat ID",
			chatID:        "",
			message:       "message",
			messageID:     1,
			expected:      0,
			err:           nil,
			errorExpected: true,
		},
		{
			name:          "error sending Telegram",
			chatID:        "12",
			message:       "message",
			messageID:     12,
			expected:      0,
			err:           errors.New("error"),
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetTelegramClient(&MockTelegram{test.expected, test.err}))
			actual, err := m.SendTelegram(context.TODO(), test.chatID, test.message)
			if test.errorExpected {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.Equal(tt, test.expected, actual)
		})
	}
}
