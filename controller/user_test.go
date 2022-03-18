package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type MockDynamoDB struct {
	getItemOutput    *dynamodb.GetItemOutput
	errGetItem       error
	putItemOutput    *dynamodb.PutItemOutput
	errPutItem       error
	deleteItemOutput *dynamodb.DeleteItemOutput
	errDeleteItem    error
	scanOutput       *dynamodb.ScanOutput
	errScan          error
}

func (m *MockDynamoDB) GetItem(
	context.Context,
	*dynamodb.GetItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return m.getItemOutput, m.errGetItem
}

func (m *MockDynamoDB) PutItem(
	context.Context,
	*dynamodb.PutItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return m.putItemOutput, m.errPutItem
}

func (m *MockDynamoDB) DeleteItem(
	context.Context,
	*dynamodb.DeleteItemInput,
	...func(*dynamodb.Options),
) (*dynamodb.DeleteItemOutput, error) {
	return m.deleteItemOutput, m.errDeleteItem
}

func (m *MockDynamoDB) Scan(
	context.Context,
	*dynamodb.ScanInput,
	...func(*dynamodb.Options),
) (*dynamodb.ScanOutput, error) {
	return m.scanOutput, m.errScan
}

func TestSuscribe(t *testing.T) {
	tests := []struct {
		name string
		chatID,
		userID,
		date int64
		kind          string
		dynamo        DynamoDBInterface
		errorExpected bool
	}{
		{
			name:          "valid suscribe",
			chatID:        12,
			userID:        10,
			date:          1647588056,
			kind:          "private",
			dynamo:        &MockDynamoDB{},
			errorExpected: false,
		},
		{
			name:          "error when suscribing",
			chatID:        12,
			userID:        10,
			date:          1647588056,
			kind:          "private",
			dynamo:        &MockDynamoDB{errPutItem: errors.New("error")},
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetDynamoDBClient(test.dynamo))
			err := m.Suscribe(context.TODO(), test.chatID, test.userID, test.date, test.kind)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestUnsuscribe(t *testing.T) {
	tests := []struct {
		name          string
		chatID        int64
		dynamo        DynamoDBInterface
		errorExpected bool
	}{
		{
			name:          "valid unsuscribe",
			chatID:        12,
			dynamo:        &MockDynamoDB{},
			errorExpected: false,
		},
		{
			name:          "error when unsuscribing",
			chatID:        12,
			dynamo:        &MockDynamoDB{errDeleteItem: errors.New("error")},
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetDynamoDBClient(test.dynamo))
			err := m.Unsuscribe(context.TODO(), test.chatID)
			if test.errorExpected {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
		})
	}
}

func TestGetChatIDs(t *testing.T) {
	tests := []struct {
		name          string
		dynamo        DynamoDBInterface
		expected      []string
		errorExpected bool
	}{
		{
			name: "valid chat IDs",
			dynamo: &MockDynamoDB{
				scanOutput: &dynamodb.ScanOutput{
					Items: []map[string]types.AttributeValue{
						{
							"ChatID": &types.AttributeValueMemberN{
								Value: "12",
							},
						},
						{
							"ChatID": &types.AttributeValueMemberN{
								Value: "13",
							},
						},
					},
				},
			},
			expected:      []string{"12", "13"},
			errorExpected: false,
		},
		{
			name:          "empty chat IDs",
			dynamo:        &MockDynamoDB{scanOutput: &dynamodb.ScanOutput{Items: []map[string]types.AttributeValue{}}},
			expected:      []string{},
			errorExpected: false,
		},
		{
			name:          "error getting chat IDs",
			dynamo:        &MockDynamoDB{errScan: errors.New("error")},
			expected:      []string{},
			errorExpected: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetDynamoDBClient(test.dynamo))
			actual, err := m.GetChatIDs(context.TODO())
			if test.errorExpected {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.ElementsMatch(tt, test.expected, actual)
		})
	}
}
