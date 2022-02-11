package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/igvaquero18/magnifibot/controller"
	"github.com/igvaquero18/magnifibot/utils"
	"github.com/mymmrac/telego"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	verboseEnv       = "MAGNIFIBOT_VERBOSE"
	awsRegionEnv     = "MAGNIFIBOT_AWS_REGION"
	sqsEndpointEnv   = "MAGNIFIBOT_SQS_ENDPOINT"
	sqsQueueNameEnv  = "MAGNIFIBOT_SQS_QUEUE_NAME"
	telegramTokenEnv = "MAGNIFIBOT_TELEGRAM_BOT_TOKEN"
)

const (
	verboseFlag       = "logging.verbose"
	awsRegionFlag     = "aws.region"
	sqsEndpointFlag   = "aws.sqs.endpoint"
	sqsQueueNameFlag  = "aws.sqs.queue_name"
	telegramTokenFlag = "telegram.bot_token"
)

const (
	telegramParseMode = "MarkdownV2"
)

var (
	c     *controller.Magnifibot
	sugar *zap.SugaredLogger
	bot   *telego.Bot
)

// Response is of type SQSEvent since we're leveraging the
// AWS SQSEvent functionality
//
// https://www.serverless.com/framework/docs/providers/aws/events/sqs
type Event events.SQSEvent

func init() {
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(awsRegionFlag, "eu-west-3")
	viper.SetDefault(sqsEndpointFlag, "")
	viper.SetDefault(sqsQueueNameFlag, controller.DefaultQueueName)
	viper.SetDefault(telegramTokenFlag, "")
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(sqsEndpointFlag, sqsEndpointEnv)
	viper.BindEnv(sqsQueueNameFlag, sqsQueueNameEnv)
	viper.BindEnv(telegramTokenFlag, telegramTokenEnv)

	var err error

	sugar, err = utils.InitSugaredLogger(viper.GetBool(verboseFlag))
	if err != nil {
		fmt.Printf("error when initializing logger: %s\n", err.Error())
		os.Exit(1)
	}

	region := viper.GetString(awsRegionFlag)
	sqsEndpoint := viper.GetString(sqsEndpointFlag)

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

	sugar.Info("creating telegram bot client")
	bot, err = telego.NewBot(viper.GetString(telegramTokenFlag), telego.WithLogger(sugar))
	if err != nil {
		sugar.Fatalw("error creating telegram bot client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetConfig(&controller.MagnifibotConfig{
			QueueURL: *queueURL.QueueUrl,
		}),
		controller.SetSQSClient(sqsClient),
	)

}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) error {
	sugar.Debug("received sqs event")

	// TODO: make this concurrent
	for _, record := range event.Records {

		// TODO: make this a method in the controller
		chatID, err := strconv.Atoi(*record.MessageAttributes["chatID"].StringValue)
		if err != nil {
			return fmt.Errorf("error converting chat ID from string to integer: %w", err)
		}

		re := regexp.MustCompile(`([_\*\[\]\(\)\~\>#\+\-\=\|\{\}\.!])`)
		gospelDay := string(re.ReplaceAll([]byte(*record.MessageAttributes["gospelDay"].StringValue), []byte(`\$1`)))
		gospelReference := string(
			re.ReplaceAll([]byte(*record.MessageAttributes["gospelReference"].StringValue), []byte(`\$1`)),
		)
		gospelTitle := string(
			re.ReplaceAll([]byte(*record.MessageAttributes["gospelTitle"].StringValue), []byte(`\$1`)),
		)
		gospelBody := string(
			re.ReplaceAll([]byte(record.Body), []byte(`\$1`)),
		)

		message, err := bot.SendMessage(&telego.SendMessageParams{
			ChatID:    telego.ChatID{ID: int64(chatID)},
			ParseMode: telegramParseMode,
			Text:      fmt.Sprintf("%s\n\n*%s\n%s*\n\n%s", gospelDay, gospelReference, gospelTitle, gospelBody),
		})
		if err != nil {
			return fmt.Errorf("error sending Telegram message: %w", err)
		}
		sugar.Debugw("successfully sent Telegram message", "chat_id", chatID, "message_id", message.MessageID)
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
