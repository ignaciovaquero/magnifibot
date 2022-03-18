package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/assert"
)

type MockQueue struct {
	output *sqs.SendMessageOutput
	err    error
}

func (m *MockQueue) SendMessage(ctx context.Context,
	params *sqs.SendMessageInput,
	optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return m.output, m.err
}

func TestSendMessageToQueue(t *testing.T) {
	tests := []struct {
		name string
		chatID,
		message string
		queue         SQSSendMessageAPI
		expected      string
		errorExpected bool
	}{
		{
			name:          "valid send message",
			chatID:        "12",
			message:       "message",
			queue:         &MockQueue{output: &sqs.SendMessageOutput{MessageId: aws.String("id")}, err: nil},
			expected:      "id",
			errorExpected: false,
		},
		{
			name:          "error sending message",
			chatID:        "12",
			message:       "message",
			queue:         &MockQueue{err: errors.New("error")},
			expected:      "",
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetSQSClient(test.queue))
			actual, err := m.SendMessageToQueue(context.TODO(), test.chatID, test.message)
			if test.errorExpected {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.Equal(tt, test.expected, actual)
		})
	}
}
