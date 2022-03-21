package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"

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

var (
	c     controller.MagnifibotInterface
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

}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, event Event) error {
	sugar.Debug("received sqs event")

	wg := sync.WaitGroup{}
	errCh := make(chan error)
	doneCh := make(chan struct{})
	wg.Add(len(event.Records))

	for _, record := range event.Records {
		go func(r events.SQSMessage, e chan<- error) {
			defer wg.Done()

			var magnificat archimadrid.Magnificat
			err := json.Unmarshal([]byte(r.Body), &magnificat)

			if err != nil {
				e <- fmt.Errorf("error unmarshalling JSON: %w", err)
				return
			}

			re := regexp.MustCompile(`([_\*\[\]\(\)\~\>#\+\-\=\|\{\}\.!])`)
			day := string(re.ReplaceAll([]byte(magnificat.Day), []byte(`\$1`)))
			firstLectureReference := string(re.ReplaceAll([]byte(magnificat.FirstLecture.Reference), []byte(`\$1`)))
			firstLectureTitle := string(re.ReplaceAll([]byte(magnificat.FirstLecture.Title), []byte(`\$1`)))
			firstLectureBody := string(re.ReplaceAll([]byte(magnificat.FirstLecture.Content), []byte(`\$1`)))
			psalmReference := string(re.ReplaceAll([]byte(magnificat.Psalm.Title), []byte(`\$1`)))
			psalmTitle := string(re.ReplaceAll([]byte(magnificat.Psalm.Reference), []byte(`\$1`)))
			psalmBody := string(re.ReplaceAll([]byte(magnificat.Psalm.Content), []byte(`\$1`)))
			gospelReference := string(re.ReplaceAll([]byte(magnificat.Gosp.Reference), []byte(`\$1`)))
			gospelTitle := string(re.ReplaceAll([]byte(magnificat.Gosp.Title), []byte(`\$1`)))
			gospelBody := string(re.ReplaceAll([]byte(magnificat.Gosp.Content), []byte(`\$1`)))

			chatID := *r.MessageAttributes["chatID"].StringValue

			// Send day
			messageID, err := c.SendTelegram(ctx, chatID, fmt.Sprintf("*%s*", day))
			if err != nil {
				e <- fmt.Errorf("error sending day %s as Telegram message: %w", day, err)
				return
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
				e <- fmt.Errorf("error sending first lecture as Telegram message: %w", err)
				return
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
			if magnificat.SecondLecture != nil {
				secondLectureReference := string(
					re.ReplaceAll([]byte(magnificat.SecondLecture.Reference), []byte(`\$1`)),
				)
				secondLectureTitle := string(re.ReplaceAll([]byte(magnificat.SecondLecture.Title), []byte(`\$1`)))
				secondLectureBody := string(re.ReplaceAll([]byte(magnificat.SecondLecture.Content), []byte(`\$1`)))

				messageID, err = c.SendTelegram(
					ctx,
					chatID,
					fmt.Sprintf("*%s\n%s*\n\n%s", secondLectureReference, secondLectureTitle, secondLectureBody),
				)
				if err != nil {
					e <- fmt.Errorf("error sending second lecture as Telegram message: %w", err)
					return
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
			messageID, err = c.SendTelegram(ctx, chatID, fmt.Sprintf("*%s\n%s*\n\n%s", psalmReference, psalmTitle, psalmBody))
			if err != nil {
				e <- fmt.Errorf("error sending psalm as Telegram message: %w", err)
				return
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
			messageID, err = c.SendTelegram(ctx, chatID, fmt.Sprintf("*%s\n%s*\n\n%s", gospelReference, gospelTitle, gospelBody))
			if err != nil {
				e <- fmt.Errorf("error sending gospel as Telegram message: %w", err)
				return
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
		}(record, errCh)
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
		return fmt.Errorf("errors occurred while processing sqs events: %s", strings.Join(errors, "\n"))
	}
	return nil
}

func main() {
	lambda.Start(Handler)
}
