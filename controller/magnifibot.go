package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Option is a function to apply settings to Magnifibot struct
type Option func(m *Magnifibot) Option

// MagnifibotInterface is the interface implemented by the SmartHome Controller
type MagnifibotInterface interface {
	Suscribe(userID, chatID, date int64, kind string) error
	Unsuscribe(chatID int64) error
}

// DynamoDBInterface is an interface implemented by the dynamodb.Client that allow
// us to mock its calls during unit testing
type DynamoDBInterface interface {
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(
		context.Context,
		*dynamodb.DeleteItemInput,
		...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
}

// Magnifibot is the controller for the Magnifibot application.
type Magnifibot struct {
	DynamoDBInterface
	Config *MagnifibotConfig
}

// MagnifibotConfig is a struct that allows to set all the configuration
// options for the Magnifibot controller
type MagnifibotConfig struct {
	// UserTable is the name of the User table in DynamoDB
	UserTable string
}

func NewMagnifibot(opts ...Option) *Magnifibot {
	m := &Magnifibot{}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// SetDynamoDBClient sets the DynamoDB client for the API
func SetDynamoDBClient(client DynamoDBInterface) Option {
	return func(m *Magnifibot) Option {
		prev := m.DynamoDBInterface
		m.DynamoDBInterface = client
		return SetDynamoDBClient(prev)
	}
}

// SetConfig sets the DynamoDB config
func SetConfig(c *MagnifibotConfig) Option {
	return func(m *Magnifibot) Option {
		prev := m.Config

		if c.UserTable == "" {
			c.UserTable = DefaultUserTable
		}

		m.Config = c
		return SetConfig(prev)
	}
}

func (m *Magnifibot) get(
	hashkey string,
	object types.AttributeValue,
	table string,
) (map[string]types.AttributeValue, error) {

	output, err := m.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key:       map[string]types.AttributeValue{hashkey: object},
		TableName: aws.String(table),
	})

	if err != nil {
		return map[string]types.AttributeValue{}, err
	}

	return output.Item, nil
}

func (m *Magnifibot) delete(hashkey string, object types.AttributeValue, table string) error {
	_, err := m.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key:       map[string]types.AttributeValue{hashkey: object},
	})
	return err
}
