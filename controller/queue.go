package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const (
	DefaultQueueName = "magnifibot"
)

// SendGospelToQueue sends the gospel for a particular ChatID to the SQS queue configured in
// the controller. It returns the Message ID on success, and an error on failure.
func (m *Magnifibot) SendMessageToQueue(ctx context.Context, chatID, message string) (string, error) {
	messageOutput, err := m.SQSSendMessageAPI.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(m.Config.QueueURL),
		MessageBody: aws.String(message),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"chatID": {DataType: aws.String("Number"), StringValue: aws.String(chatID)},
		},
	})
	if err != nil {
		return "", err
	}
	return *messageOutput.MessageId, nil
}
