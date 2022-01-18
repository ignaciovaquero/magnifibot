package controller

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/igvaquero18/magnifibot/utils"
)

// Option is a function to apply settings to Magnifibot struct
type Option func(m *Magnifibot) Option

// MagnifibotInterface is the interface implemented by the SmartHome Controller
type MagnifibotInterface interface {
}

// DynamoDBInterface is an interface implemented by the dynamodb.Client that allow
// us to mock its calls during unit testing
type DynamoDBInterface interface {
	GetItem(context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	DeleteItem(context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// Magnifibot is the controller for the Magnifibot application.
type Magnifibot struct {
	utils.Logger
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
