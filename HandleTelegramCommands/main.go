package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igvaquero18/magnifibot/api"
	"github.com/igvaquero18/magnifibot/controller"
	"github.com/igvaquero18/magnifibot/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	magnifibotNameEnv    = "MAGNIFIBOT_NAME"
	awsRegionEnv         = "MAGNIFIBOT_AWS_REGION"
	verboseEnv           = "MAGNIFIBOT_VERBOSE"
	dynamoDBEndpointEnv  = "MAGNIFIBOT_DYNAMODB_ENDPOINT"
	dynamoDBUserTableEnv = "MAGNIFIBOT_DYNAMODB_USER_TABLE"
	magnifibotTimeoutEnv = "MAGNIFIBOT_TIMEOUT"
)

const (
	magnifibotNameFlag    = "name"
	awsRegionFlag         = "aws.region"
	verboseFlag           = "logging.verbose"
	dynamoDBEndpointFlag  = "aws.dynamodb.endpoint"
	dynamoDBUserTableFlag = "aws.dynamodb.tables.user"
	magnifibotTimeoutFlag = "timeout"
)

var (
	c     *controller.Magnifibot
	sugar *zap.SugaredLogger
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

func init() {
	viper.SetDefault(magnifibotNameFlag, "magnifibot_bot")
	viper.SetDefault(awsRegionFlag, "eu-west-3")
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBUserTableFlag, controller.DefaultUserTable)
	viper.SetDefault(magnifibotTimeoutFlag, utils.DefaultTimeout)
	viper.BindEnv(magnifibotNameFlag, magnifibotNameEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBUserTableFlag, dynamoDBUserTableEnv)
	viper.BindEnv(magnifibotTimeoutFlag, magnifibotTimeoutEnv)

	var err error

	sugar, err = utils.InitSugaredLogger(viper.GetBool(verboseFlag))
	if err != nil {
		fmt.Printf("error when initializing logger: %s\n", err.Error())
		os.Exit(1)
	}

	region := viper.GetString(awsRegionFlag)
	dynamoDBEndpoint := viper.GetString(dynamoDBEndpointFlag)

	sugar.Infow("creating DynamoDB client", "region", region, "url", dynamoDBEndpoint)
	dynamoClient, err := utils.InitDynamoClient(region, dynamoDBEndpoint)
	if err != nil {
		sugar.Fatalw("error creating DynamoDB client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetDynamoDBClient(dynamoClient),
		controller.SetConfig(&controller.MagnifibotConfig{
			UserTable: viper.GetString(dynamoDBUserTableFlag),
		}),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	ctx, cancel, err := utils.InitContextWithTimeout(viper.GetString(magnifibotTimeoutFlag))
	if err != nil {
		sugar.Warnw("invalid timeout setting", "timeout", viper.GetString(magnifibotTimeoutFlag), "default", "3s", "error", err.Error())
	}
	defer cancel()
	sugar.Infow("received request", "method", request.HTTPMethod, "body", request.Body)
	var update api.Update
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	if err := json.Unmarshal([]byte(request.Body), &update); err != nil {
		return Response{
			Body:       fmt.Sprintf("Invalid Telegram message: %s", err.Error()),
			StatusCode: http.StatusBadRequest,
			Headers:    headers,
		}, nil
	}

	if update.Message != nil {
		sugar.Infow(
			"received message",
			"chat_id",
			update.Message.Chat.ID,
			"text",
			update.Message.Text,
		)
		re := regexp.MustCompile(`/(\w*)@?\w*`)
		submatch := re.FindStringSubmatch(update.Message.Text)
		if submatch != nil {
			command := api.ToCommand(submatch[1])
			if !command.IsValid() {
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.Message.Chat.ID,
					fmt.Sprintf(
						"Lo siento, solo acepto los siguientes comandos: %s",
						strings.Join(api.GetValidCommandsString(), ", "),
					),
				)
			}
			if api.ValidCommands["suscribe"] == command {
				sugar.Infow("suscribing user", "chat_id", update.Message.Chat.ID)
				if err := c.Suscribe(ctx, update.Message.Chat.ID, update.Message.From.ID, update.Message.Date, update.Message.Chat.Type); err != nil {
					return createTelegramResponse(
						http.StatusOK,
						headers,
						update.Message.Chat.ID,
						fmt.Sprintf("Lo siento, no pude suscribirte: %s", err.Error()),
					)
				}
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.Message.Chat.ID,
					"¡Hecho! Te enviaré el Evangelio cada día.",
				)
			}
			sugar.Infow("unsuscribing user", "chat_id", update.Message.Chat.ID)
			if err := c.Unsuscribe(ctx, update.Message.Chat.ID); err != nil {
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.Message.Chat.ID,
					fmt.Sprintf("Lo siento, no pude darte de baja: %s", err.Error()),
				)
			}
			return createTelegramResponse(
				http.StatusOK,
				headers,
				update.Message.Chat.ID,
				"¡Hecho! Ya no te enviaré más el Evangelio.",
			)
		}
		return createTelegramResponse(
			http.StatusOK,
			headers,
			update.Message.Chat.ID,
			"Lo siento, solo acepto comandos de Telegram.",
		)
	}

	if update.ChannelPost != nil {
		sugar.Infow(
			"received channel post",
			"chat_id",
			update.ChannelPost.Chat.ID,
			"text",
			update.ChannelPost.Text,
		)
		re := regexp.MustCompile(fmt.Sprintf(`/(\w*)@%s`, viper.GetString(magnifibotNameFlag)))
		submatch := re.FindStringSubmatch(update.ChannelPost.Text)
		if submatch != nil {
			command := api.ToCommand(submatch[1])
			if !command.IsValid() {
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.ChannelPost.Chat.ID,
					fmt.Sprintf(
						"Lo siento, solo acepto los siguientes comandos: %s",
						strings.Join(api.GetValidCommandsString(), ", "),
					),
				)
			}
			if api.ValidCommands["suscribe"] == command {
				sugar.Infow("suscribing user", "chat_id", update.ChannelPost.Chat.ID)
				if err := c.Suscribe(ctx, update.ChannelPost.Chat.ID, update.ChannelPost.SenderChat.ID, update.ChannelPost.Date, update.ChannelPost.SenderChat.Type); err != nil {
					return createTelegramResponse(
						http.StatusOK,
						headers,
						update.ChannelPost.Chat.ID,
						fmt.Sprintf("Lo siento, no pude suscribirte: %s", err.Error()),
					)
				}
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.ChannelPost.Chat.ID,
					"¡Hecho! Te enviaré el Evangelio cada día.",
				)
			}
			sugar.Infow("unsuscribing user", "chat_id", update.ChannelPost.Chat.ID)
			if err := c.Unsuscribe(ctx, update.ChannelPost.Chat.ID); err != nil {
				return createTelegramResponse(
					http.StatusOK,
					headers,
					update.ChannelPost.Chat.ID,
					fmt.Sprintf("Lo siento, no pude darte de baja: %s", err.Error()),
				)
			}
			return createTelegramResponse(
				http.StatusOK,
				headers,
				update.ChannelPost.Chat.ID,
				"¡Hecho! Ya no te enviaré más el Evangelio.",
			)
		}
		return createTelegramResponse(
			http.StatusOK,
			headers,
			update.ChannelPost.Chat.ID,
			"Lo siento, solo acepto comandos de Telegram.",
		)
	}
	return Response{
		StatusCode: http.StatusBadRequest,
		Headers:    headers,
		Body:       "Invalid request",
	}, nil
}

func createTelegramResponse(
	status int,
	headers map[string]string,
	chatID int64,
	text string,
) (Response, error) {
	t := api.TelegramWebhookSendMessage{
		ChatID: chatID,
		Text:   text,
		Method: "sendMessage",
	}

	body, err := json.Marshal(t)
	if err != nil {
		return Response{
			Body:       fmt.Sprintf("error when marshalling response: %s", err.Error()),
			StatusCode: http.StatusInternalServerError,
			Headers:    headers,
		}, fmt.Errorf("error when marshalling response: %w", err)
	}

	return Response{
		Body:       string(body),
		StatusCode: status,
		Headers:    headers,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
