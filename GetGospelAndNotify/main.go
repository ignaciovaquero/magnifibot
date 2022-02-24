package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/igvaquero18/magnifibot/archimadrid"
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
	c      *controller.Magnifibot
	sugar  *zap.SugaredLogger
	gospel *archimadrid.Gospel
)

// Response is of type CloudWatchEvent since we're leveraging the
// AWS CloudWatch Event functionality
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Event events.CloudWatchEvent

func init() {
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(awsRegionFlag, "eu-west-3")
	viper.SetDefault(sqsEndpointFlag, "")
	viper.SetDefault(sqsQueueNameFlag, controller.DefaultQueueName)
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

	sugar.Infow("creating SQS client", "region", region, "url", viper.GetString(sqsEndpointFlag))
	sqsClient, err := utils.InitSQSClient(region, sqsEndpoint)
	if err != nil {
		sugar.Fatalw("error creating SQS client", "error", err.Error())
	}

	queueURL, err := sqsClient.GetQueueUrl(context.TODO(), &sqs.GetQueueUrlInput{
		QueueName: aws.String(viper.GetString(sqsQueueNameFlag)),
	})

	if err != nil {
		sugar.Fatalw(
			"error getting the queue URL",
			"queue_name",
			viper.GetString(sqsQueueNameFlag),
			"error",
			err.Error(),
		)
	}

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetConfig(&controller.MagnifibotConfig{
			UserTable: viper.GetString(dynamoDBUserTableFlag),
			QueueURL:  *queueURL.QueueUrl,
		}),
		controller.SetSQSClient(sqsClient),
		controller.SetDynamoDBClient(dynamoClient),
	)

	sugar.Infow("getting gospel for day", "day", time.Now().Format("2006-01-02"))
	if gospel, err = archimadrid.NewClient().GetGospel(time.Now()); err != nil {
		sugar.Fatalw("error getting gospel", "error", err.Error())
	}
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) error {
	sugar.Infow("received cloudwatch event", "time", event.Time)
	chatIDs, err := c.GetChatIDs(ctx)
	if err != nil {
		sugar.Errorw("error getting chat IDs", "error", err.Error())
		return err
	}

	wg := &sync.WaitGroup{}
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg.Add(len(chatIDs))

	for _, chatID := range chatIDs {
		go func(e chan<- error, id string) {
			sugar.Debugw("sending message to queue", "queue_url", c.Config.QueueURL, "chat_id", id)

			messageID, err := c.SendGospelToQueue(ctx, id, gospel)

			if err != nil {
				e <- fmt.Errorf("error sending message to queue: %w", err)
				wg.Done()
				return
			}

			sugar.Debugw(
				"message stored in SQS queue",
				"queue_url",
				viper.GetString(c.Config.QueueURL),
				"message_id",
				messageID,
			)

			wg.Done()
		}(errCh, chatID)
	}

	go func(d chan<- struct{}) {
		wg.Wait()
		close(d)
	}(doneCh)

	done := false
	errors := []string{}

	for !done {
		select {
		case err := <-errCh:
			errors = append(errors, err.Error())
		case <-doneCh:
			done = true
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors while sending messages to queue: %v", strings.Join(errors, "\n"))
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
