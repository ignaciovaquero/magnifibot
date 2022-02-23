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

func (m *Magnifibot) Suscribe(ctx context.Context, chatID, userID, date int64, kind string) error {
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

func (m *Magnifibot) GetChatIDs(ctx context.Context) ([]string, error) {
	scanOutput, err := m.Scan(ctx, &dynamodb.ScanInput{
		TableName:            aws.String(m.Config.UserTable),
		ProjectionExpression: aws.String("ChatID"),
	})
	if err != nil {
		return []string{}, fmt.Errorf("error scanning dynamodb table: %w", err)
	}

	chatIDs := []string{}
	for _, item := range scanOutput.Items {
		if chatID, ok := item["ChatID"]; ok {
			id, ok := chatID.(*types.AttributeValueMemberN)
			if !ok {
				return []string{}, fmt.Errorf("error converting ChatID into a string: %w", err)
			}
			chatIDs = append(chatIDs, id.Value)
		} else {
			return []string{}, fmt.Errorf("error getting ChatID for item %v: %w", item, err)
		}
	}

	return chatIDs, nil
}
