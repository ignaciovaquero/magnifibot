package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igvaquero18/magnifibot/controller"
	"github.com/igvaquero18/magnifibot/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	verboseEnv     = "MAGNIFIBOT_VERBOSE"
	awsRegionEnv   = "MAGNIFIBOT_AWS_REGION"
	sqsNameEnv     = "MAGNIFIBOT_SQS_QUEUE_NAME"
	sqsEndpointEnv = "MAGNIFIBOT_SQS_ENDPOINT"
)

const (
	verboseFlag     = "logging.verbose"
	awsRegionFlag   = "aws.region"
	sqsNameFlag     = "aws.sqs.name"
	sqsEndpointFlag = "aws.sqs.endpoint"
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
	viper.SetDefault(sqsNameFlag, "magnifibot")
	viper.SetDefault(sqsEndpointEnv, "")
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(sqsNameFlag, sqsNameEnv)
	viper.BindEnv(sqsEndpointFlag, sqsEndpointEnv)

	var err error

	sugar, err = utils.InitSugaredLogger(viper.GetBool(verboseFlag))
	if err != nil {
		fmt.Printf("error when initializing logger: %s\n", err.Error())
		os.Exit(1)
	}

	region := viper.GetString(awsRegionFlag)
	sqsEndpoint := viper.GetString(sqsEndpointFlag)

	sugar.Infow("creating SQS client", "region", region, "url", sqsEndpointFlag)
	sqsClient, err := utils.InitSQSClient(region, sqsEndpoint)
	if err != nil {
		sugar.Fatalw("error creating SQS client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetSQSClient(sqsClient),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) (string, error) {
	sugar.Infow("received cloudwatch event", "time", event.Time)
	return "", nil
}

func main() {
	lambda.Start(Handler)
}
