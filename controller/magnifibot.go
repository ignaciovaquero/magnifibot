package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/igvaquero18/magnifibot/archimadrid"
	"github.com/mymmrac/telego"
)

// Option is a function to apply settings to Magnifibot struct
type Option func(m *Magnifibot) Option

// MagnifibotInterface is the interface implemented by the SmartHome Controller
type MagnifibotInterface interface {
	Suscribe(ctx context.Context, chatID, userID, date int64, kind string) error
	Unsuscribe(ctx context.Context, chatID int64) error
	GetChatIDs(ctx context.Context) ([]string, error)
	SendGospelToQueue(ctx context.Context, chatID string, gospel *archimadrid.Gospel) (string, error)
	GetConfig() *MagnifibotConfig
	SendTelegram(ctx context.Context, chatID string, message string) (int, error)
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
	dynamodb.ScanAPIClient
}

// SQSSendMessageAPI defines the interface for the GetQueueUrl and SendMessage functions.
// We use this interface to test the functions using a mocked service.
type SQSSendMessageAPI interface {
	SendMessage(ctx context.Context,
		params *sqs.SendMessageInput,
		optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// TelegramAPI is the interface implemented by the Telegram API
type TelegramAPI interface {
	SendMessage(params *telego.SendMessageParams) (*telego.Message, error)
}

// Magnifibot is the controller for the Magnifibot application.
type Magnifibot struct {
	DynamoDBInterface
	SQSSendMessageAPI
	TelegramAPI
	Config *MagnifibotConfig
}

// MagnifibotConfig is a struct that allows to set all the configuration
// options for the Magnifibot controller
type MagnifibotConfig struct {
	// UserTable is the name of the User table in DynamoDB
	UserTable string

	// QueueURL is the URL of the SQS queue
	QueueURL string
}

func NewMagnifibot(opts ...Option) *Magnifibot {
	m := &Magnifibot{
		Config: &MagnifibotConfig{
			UserTable: DefaultUserTable,
		},
	}

	for _, opt := range opts {
		opt(m)
	}
	return m
}

// SetDynamoDBClient sets the DynamoDB client
func SetDynamoDBClient(client DynamoDBInterface) Option {
	return func(m *Magnifibot) Option {
		prev := m.DynamoDBInterface
		m.DynamoDBInterface = client
		return SetDynamoDBClient(prev)
	}
}

// SetSQSClient sets the SQS client
func SetSQSClient(client SQSSendMessageAPI) Option {
	return func(m *Magnifibot) Option {
		prev := m.SQSSendMessageAPI
		m.SQSSendMessageAPI = client
		return SetSQSClient(prev)
	}
}

// SetTelegramClient sets the Telegram client
func SetTelegramClient(client TelegramAPI) Option {
	return func(m *Magnifibot) Option {
		prev := m.TelegramAPI
		m.TelegramAPI = client
		return SetTelegramClient(prev)
	}
}

// SetConfig sets the Magnifibot config
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

// GetConfig gets the Magnifibot config
func (m *Magnifibot) GetConfig() *MagnifibotConfig {
	return m.Config
}

func (m *Magnifibot) delete(ctx context.Context, hashkey string, object types.AttributeValue, table string) error {
	_, err := m.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(table),
		Key:       map[string]types.AttributeValue{hashkey: object},
	})
	return err
}
