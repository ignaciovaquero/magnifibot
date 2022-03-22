package main

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/igvaquero18/magnifibot/archimadrid"
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
	sqsQueueNameEnv  = "MAGNIFIBOT_SQS_ON_DEMAND_QUEUE_NAME"
	telegramTokenEnv = "MAGNIFIBOT_TELEGRAM_BOT_TOKEN"
)

const (
	verboseFlag       = "logging.verbose"
	awsRegionFlag     = "aws.region"
	sqsEndpointFlag   = "aws.sqs.endpoint"
	sqsQueueNameFlag  = "aws.sqs.on_demand_queue_name"
	telegramTokenFlag = "telegram.bot_token"
)

var (
	c     controller.MagnifibotInterface
	a     archimadrid.Archimadrid
	sugar *zap.SugaredLogger
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
	viper.SetDefault(sqsQueueNameFlag, controller.DefaultOnDemandQueueName)
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
	bot, err := telego.NewBot(viper.GetString(telegramTokenFlag), telego.WithLogger(sugar))
	if err != nil {
		sugar.Fatalw("error creating telegram bot client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetConfig(&controller.MagnifibotConfig{
			QueueURL: *queueURL.QueueUrl,
		}),
		controller.SetSQSClient(sqsClient),
		controller.SetTelegramClient(bot),
	)

	a = archimadrid.NewClient()
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) error {
	sugar.Debug("received sqs event")
	today := time.Now()

	for _, record := range event.Records {
		gospel, err := a.GetGospel(ctx, today)
		if err != nil {
			sugar.Fatalw("error getting gospel", "error", err.Error())
		}
		firstLecture, err := a.GetFirstLecture(ctx, today)
		if err != nil {
			sugar.Fatalw("error getting first lecture", "error", err.Error())
		}
		psalm, err := a.GetPsalm(ctx, today)
		if err != nil {
			sugar.Fatalw("error getting psalm", "error", err.Error())
		}
		secondLecture, err := a.GetSecondLecture(ctx, today)
		if err != nil {
			sugar.Fatalw("error getting second lecture", "error", err.Error())
		}

		re := regexp.MustCompile(`([_\*\[\]\(\)\~\>#\+\-\=\|\{\}\.!])`)
		day := string(re.ReplaceAll([]byte(gospel.Day), []byte(`\$1`)))
		firstLectureReference := string(re.ReplaceAll([]byte(firstLecture.Reference), []byte(`\$1`)))
		firstLectureTitle := string(re.ReplaceAll([]byte(firstLecture.Title), []byte(`\$1`)))
		firstLectureBody := string(re.ReplaceAll([]byte(firstLecture.Content), []byte(`\$1`)))
		psalmReference := string(re.ReplaceAll([]byte(psalm.Title), []byte(`\$1`)))
		psalmTitle := string(re.ReplaceAll([]byte(psalm.Reference), []byte(`\$1`)))
		psalmBody := string(re.ReplaceAll([]byte(psalm.Content), []byte(`\$1`)))
		gospelReference := string(re.ReplaceAll([]byte(gospel.Reference), []byte(`\$1`)))
		gospelTitle := string(re.ReplaceAll([]byte(gospel.Title), []byte(`\$1`)))
		gospelBody := string(re.ReplaceAll([]byte(gospel.Content), []byte(`\$1`)))

		chatID := *record.MessageAttributes["chatID"].StringValue

		// Send day
		messageID, err := c.SendTelegram(ctx, chatID, fmt.Sprintf("*%s*", day))
		if err != nil {
			return fmt.Errorf("error sending day %s as Telegram message: %w", day, err)
		}
		sugar.Debugw(
			"successfully sent day as Telegram message",
			"day",
			day,
			"chat_id",
			chatID,
			"message_id",
			messageID,
		)

		// Send First Lecture
		messageID, err = c.SendTelegram(
			ctx,
			chatID,
			fmt.Sprintf("*%s\n%s*\n\n%s", firstLectureReference, firstLectureTitle, firstLectureBody),
		)
		if err != nil {
			return fmt.Errorf("error sending first lecture as Telegram message: %w", err)
		}
		sugar.Debugw(
			"successfully sent first lecture as Telegram message",
			"first_lecture_reference",
			firstLectureReference,
			"chat_id",
			chatID,
			"message_id",
			messageID,
		)

		// Send Second Lecture if exists
		if len(secondLecture.Content) > 0 {
			secondLectureReference := string(
				re.ReplaceAll([]byte(secondLecture.Reference), []byte(`\$1`)),
			)
			secondLectureTitle := string(re.ReplaceAll([]byte(secondLecture.Title), []byte(`\$1`)))
			secondLectureBody := string(re.ReplaceAll([]byte(secondLecture.Content), []byte(`\$1`)))

			messageID, err = c.SendTelegram(
				ctx,
				chatID,
				fmt.Sprintf("*%s\n%s*\n\n%s", secondLectureReference, secondLectureTitle, secondLectureBody),
			)
			if err != nil {
				return fmt.Errorf("error sending second lecture as Telegram message: %w", err)
			}
			sugar.Debugw(
				"successfully sent second lecture as Telegram message",
				"second_lecture_reference",
				secondLectureReference,
				"chat_id",
				chatID,
				"message_id",
				messageID,
			)
		}

		// Send Psalm
		messageID, err = c.SendTelegram(
			ctx,
			chatID,
			fmt.Sprintf("*%s\n%s*\n\n%s", psalmReference, psalmTitle, psalmBody),
		)
		if err != nil {
			return fmt.Errorf("error sending psalm as Telegram message: %w", err)
		}
		sugar.Debugw(
			"successfully sent psalm as Telegram message",
			"psalm_reference",
			psalmTitle,
			"chat_id",
			chatID,
			"message_id",
			messageID,
		)

		// Send Gospel
		messageID, err = c.SendTelegram(
			ctx,
			chatID,
			fmt.Sprintf("*%s\n%s*\n\n%s", gospelReference, gospelTitle, gospelBody),
		)
		if err != nil {
			return fmt.Errorf("error sending gospel as Telegram message: %w", err)
		}
		sugar.Debugw(
			"successfully sent gospel as Telegram message",
			"gospel_reference",
			gospelReference,
			"chat_id",
			chatID,
			"message_id",
			messageID,
		)
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
