package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
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
	sqsEndpointEnv       = "MAGNIFIBOT_SQS_ENDPOINT"
	sqsQueueNameEnv      = "MAGNIFIBOT_SQS_QUEUE_NAME"
	dynamoDBEndpointEnv  = "MAGNIFIBOT_DYNAMODB_ENDPOINT"
	dynamoDBUserTableEnv = "MAGNIFIBOT_DYNAMODB_USER_TABLE"
)

const (
	verboseFlag           = "logging.verbose"
	awsRegionFlag         = "aws.region"
	sqsEndpointFlag       = "aws.sqs.endpoint"
	sqsQueueNameFlag      = "aws.sqs.queue_name"
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
	viper.SetDefault(sqsEndpointEnv, "")
	viper.SetDefault(sqsQueueNameEnv, controller.DefaultQueueName)
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBUserTableFlag, controller.DefaultUserTable)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(sqsEndpointFlag, sqsEndpointEnv)
	viper.BindEnv(sqsQueueNameFlag, sqsQueueNameEnv)
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
		controller.SetConfig(&controller.MagnifibotConfig{
			UserTable: viper.GetString(dynamoDBUserTableFlag),
			QueueURL:  viper.GetString(sqsEndpointFlag),
		}),
		controller.SetSQSClient(sqsClient),
		controller.SetDynamoDBClient(dynamoClient),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) (string, error) {
	sugar.Infow("received cloudwatch event", "time", event.Time)
	scanOutput, err := c.Scan(ctx, &dynamodb.ScanInput{
		TableName:            aws.String(c.Config.UserTable),
		ProjectionExpression: aws.String("ChatID"),
	})
	if err != nil {
		return "", fmt.Errorf("error scanning dynamodb table: %w", err)
	}
	// TODO: Make concurrent
	for _, item := range scanOutput.Items {
		if chatID, ok := item["ChatID"]; ok {
			sugar.Debugw("sending message to chat", "chat_id", chatID)
			messageOutput, err := c.SendMessage(ctx, &sqs.SendMessageInput{
				QueueUrl:    aws.String(c.Config.QueueURL),
				MessageBody: aws.String(fmt.Sprintf("ChatID: %s", chatID)),
			})
			if err != nil {
				return "", fmt.Errorf("error sending message to queue: %w", err)
			}
			sugar.Debugw(
				"message stored in SQS queue",
				"queue_url",
				viper.GetString(sqsQueueURLFlag),
				"message_id",
				*messageOutput.MessageId,
			)
		} else {
			sugar.Errorw("error getting ChatID", "item", item)
		}
	}
	return "", nil
}

// func main() {
// 	lambda.Start(Handler)
// }

func main() {
	result, err := Handler(context.TODO(), Event{Time: time.Now()})
	if err != nil {
		sugar.Fatalw("error handling event", "error", err.Error())
	}
	fmt.Println(result)
}
