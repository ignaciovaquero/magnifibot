package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/igvaquero18/magnifibot/controller"
	"github.com/igvaquero18/magnifibot/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	verboseEnv           = "MAGNIFIBOT_VERBOSE"
	awsRegionEnv         = "MAGNIFIBOT_AWS_REGION"
	sqsNameEnv           = "MAGNIFIBOT_SQS_QUEUE_NAME"
	sqsEndpointEnv       = "MAGNIFIBOT_SQS_ENDPOINT"
	dynamoDBEndpointEnv  = "MAGNIFIBOT_DYNAMODB_ENDPOINT"
	dynamoDBUserTableEnv = "MAGNIFIBOT_DYNAMODB_USER_TABLE"
)

const (
	verboseFlag           = "logging.verbose"
	awsRegionFlag         = "aws.region"
	sqsNameFlag           = "aws.sqs.name"
	sqsEndpointFlag       = "aws.sqs.endpoint"
	dynamoDBEndpointFlag  = "aws.dynamodb.endpoint"
	dynamoDBUserTableFlag = "aws.dynamodb.tables.user"
)

var (
	c     *controller.Magnifibot
	sugar *zap.SugaredLogger
)

// Response is of type CloudWatchEvent since we're leveraging the
// AWS CloudWatch Event functionality
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Event events.CloudWatchEvent

func init() {
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(awsRegionFlag, "eu-west-3")
	viper.SetDefault(sqsNameFlag, "magnifibot")
	viper.SetDefault(sqsEndpointEnv, "")
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBUserTableFlag, controller.DefaultUserTable)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(sqsNameFlag, sqsNameEnv)
	viper.BindEnv(sqsEndpointFlag, sqsEndpointEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBUserTableFlag, dynamoDBUserTableEnv)

	var err error

	sugar, err = utils.InitSugaredLogger(viper.GetBool(verboseFlag))
	if err != nil {
		fmt.Printf("error when initializing logger: %s\n", err.Error())
		os.Exit(1)
	}

	region := viper.GetString(awsRegionFlag)
	sqsEndpoint := viper.GetString(sqsEndpointFlag)
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)

	sugar.Infow("creating SQS client", "region", region, "url", sqsEndpointFlag)
	sqsClient, err := utils.InitSQSClient(region, sqsEndpoint)
	if err != nil {
		sugar.Fatalw("error creating SQS client", "error", err.Error())
	}

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetSQSClient(sqsClient),
		controller.SetDynamoDBClient(dynamoClient),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) (string, error) {
	sugar.Infow("received cloudwatch event", "time", event.Time)
	scanOutput, err := c.Scan(ctx, &dynamodb.ScanInput{
		ProjectionExpression: aws.String("ChatID"),
	})
	if err != nil {
		return "", fmt.Errorf("error scanning dynamodb table: %w", err)
	}
	// TODO: Make concurrent
	for _, item := range scanOutput.Items {
		if chatID, ok := item["ChatID"]; ok {
			sugar.Infow("sending message to chat", "chat_id", chatID)
			c.SendMessage(ctx, &sqs.SendMessageInput{
				MessageBody: aws.String(fmt.Sprintf("ChatID: %s", chatID)),
			})
		} else {
			sugar.Errorw("error getting ChatID", "item", item)
		}
	}
	return "", nil
}

func main() {
	lambda.Start(Handler)
}
