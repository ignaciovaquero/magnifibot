package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/igvaquero18/magnifibot/archimadrid"
)

const (
	DefaultQueueName = "magnifibot"
)

// SendGospelToQueue sends the gospel for a particular ChatID to the SQS queue configured in
// the controller. It returns the Message ID on success, and an error on failure.
func (m *Magnifibot) SendGospelToQueue(ctx context.Context, chatID string, gospel *archimadrid.Gospel) (string, error) {
	messageOutput, err := m.SQSSendMessageAPI.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(m.Config.QueueURL),
		MessageBody: aws.String(gospel.Content),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"chatID":          {DataType: aws.String("Number"), StringValue: aws.String(chatID)},
			"gospelTitle":     {DataType: aws.String("String"), StringValue: aws.String(gospel.Title)},
			"gospelDay":       {DataType: aws.String("String"), StringValue: aws.String(gospel.Day)},
			"gospelReference": {DataType: aws.String("String"), StringValue: aws.String(gospel.Reference)},
		},
	})
	if err != nil {
		return "", err
	}
	return *messageOutput.MessageId, nil
}
