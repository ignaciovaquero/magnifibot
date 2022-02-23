package controller

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	DefaultUserTable = "MagnifibotUser"
)

func (m *Magnifibot) Suscribe(ctx context.Context, userID, date int64, kind string) error {
	_, err := m.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(m.Config.UserTable),
		Item: map[string]types.AttributeValue{
			"ChatID": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", chatID)},
			"ID":     &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", userID)},
			"Date":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", date)},
			"Kind":   &types.AttributeValueMemberS{Value: kind},
		},
	})

	if err != nil {
		return fmt.Errorf("error when suscribing user %d: %w", userID, err)
	}
	return nil
}

func (m *Magnifibot) Unsuscribe(ctx context.Context, chatID int64) error {
	if err := m.delete(ctx, "ChatID", &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", chatID)}, m.Config.UserTable); err != nil {
		return fmt.Errorf("error when deleting chat with id %d: %w", chatID, err)
	}
	return nil
}

func (m *Magnifibot) GetChats(ctx context.Context) ([]int64, error) {
	scanOutput, err := m.Scan(ctx, &dynamodb.ScanInput{
		TableName:            aws.String(m.Config.UserTable),
		ProjectionExpression: aws.String("ChatID"),
	})

	return []int64{}, nil
}
