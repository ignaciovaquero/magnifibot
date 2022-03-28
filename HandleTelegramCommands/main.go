package main

import (
	"context"
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
	sqsEndpointEnv       = "MAGNIFIBOT_SQS_ENDPOINT"
	verboseEnv           = "MAGNIFIBOT_VERBOSE"
	dynamoDBEndpointEnv  = "MAGNIFIBOT_DYNAMODB_ENDPOINT"
	dynamoDBUserTableEnv = "MAGNIFIBOT_DYNAMODB_USER_TABLE"
	lambdaEndpointEnv    = "MAGNIFIBOT_LAMBDA_ENDPOINT"
	onDemandLambdaEnv    = "MAGNIFIBOT_ON_DEMAND_LAMBDA_FUNCTION_NAME"
	magnifibotTimeoutEnv = "MAGNIFIBOT_TIMEOUT"
)

const (
	magnifibotNameFlag    = "name"
	awsRegionFlag         = "aws.region"
	sqsEndpointFlag       = "aws.sqs.endpoint"
	verboseFlag           = "logging.verbose"
	dynamoDBEndpointFlag  = "aws.dynamodb.endpoint"
	dynamoDBUserTableFlag = "aws.dynamodb.tables.user"
	lambdaEndpointFlag    = "aws.lambda.endpoint"
	onDemandLambdaFlag    = "aws.lambda.on_demand.function_name"
	magnifibotTimeoutFlag = "timeout"
)

var (
	c     controller.MagnifibotInterface
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
	viper.SetDefault(sqsEndpointFlag, "")
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBUserTableFlag, controller.DefaultUserTable)
	viper.SetDefault(lambdaEndpointFlag, "")
	viper.SetDefault(onDemandLambdaFlag, "")
	viper.SetDefault(magnifibotTimeoutFlag, utils.DefaultTimeout)
	viper.BindEnv(magnifibotNameFlag, magnifibotNameEnv)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(sqsEndpointFlag, sqsEndpointEnv)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBUserTableFlag, dynamoDBUserTableEnv)
	viper.BindEnv(lambdaEndpointFlag, lambdaEndpointEnv)
	viper.BindEnv(onDemandLambdaFlag, onDemandLambdaEnv)
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

	lambdaEndpoint := viper.GetString(lambdaEndpointFlag)
	sugar.Infow("creating lambda client", "region", region, "url", lambdaEndpoint)
	lambdaClient, err := utils.InitLambdaClient(region, lambdaEndpoint)
	if err != nil {
		sugar.Fatalw("error creating Lambda client", "error", err.Error())
	}

	c = controller.NewMagnifibot(
		controller.SetDynamoDBClient(dynamoClient),
		controller.SetLambdaClient(lambdaClient),
		controller.SetConfig(&controller.MagnifibotConfig{
			UserTable: viper.GetString(dynamoDBUserTableFlag),
		}),
	)
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	ctx, cancel, err := utils.InitContextWithTimeout(viper.GetString(magnifibotTimeoutFlag))
	if err != nil {
		sugar.Warnw(
			"invalid timeout setting, using default",
			"timeout",
			viper.GetString(magnifibotTimeoutFlag),
			"default",
			utils.DefaultTimeout,
			"error",
			err.Error(),
		)
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
		return handleCommand(
			ctx,
			`/(\w*)@?\w*`,
			update.Message.Text,
			update.Message.Chat.Type,
			update.Message.Chat.ID,
			update.Message.From.ID,
			update.Message.Date,
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
		return handleCommand(
			ctx,
			fmt.Sprintf(`/(\w*)@%s`, viper.GetString(magnifibotNameFlag)),
			update.ChannelPost.Text,
			update.ChannelPost.SenderChat.Type,
			update.ChannelPost.Chat.ID,
			update.ChannelPost.SenderChat.ID,
			update.ChannelPost.Date,
		)
	}
	return Response{
		StatusCode: http.StatusBadRequest,
		Headers:    headers,
		Body:       "Invalid request",
	}, nil
}

func handleCommand(ctx context.Context, regexPattern, body, kind string, chatID, userID, date int64) (Response, error) {
	re := regexp.MustCompile(regexPattern)
	submatch := re.FindStringSubmatch(body)
	if submatch != nil {
		command := api.ToCommand(submatch[1])
		if !command.IsValid() {
			return createTelegramResponse(
				http.StatusOK,
				chatID,
				fmt.Sprintf(
					"Lo siento, solo acepto los siguientes comandos: %s",
					strings.Join(api.GetValidCommandsString(), ", "),
				),
			)
		}
		if api.ValidCommands["suscribe"] == command {
			sugar.Infow("suscribe operation", "chat_id", chatID)
			if err := c.Suscribe(ctx, chatID, userID, date, kind); err != nil {
				return createTelegramResponse(
					http.StatusOK,
					chatID,
					fmt.Sprintf("Lo siento, no he podido suscribirte: %s", err.Error()),
				)
			}
			return createTelegramResponse(
				http.StatusOK,
				chatID,
				"¡Hecho! Te enviaré el Evangelio cada día.",
			)
		} else if api.ValidCommands["unsuscribe"] == command {
			sugar.Infow("unsuscribe operation", "chat_id", chatID)
			if err := c.Unsuscribe(ctx, chatID); err != nil {
				return createTelegramResponse(
					http.StatusOK,
					chatID,
					fmt.Sprintf("Lo siento, no he podido darte de baja: %s", err.Error()),
				)
			}
			return createTelegramResponse(
				http.StatusOK,
				chatID,
				"¡Hecho! Ya no te enviaré más el Evangelio.",
			)
		}

		sugar.Infow("on demand operation", "chat_id", chatID)
		lambdaFunctionName := viper.GetString(onDemandLambdaFlag)
		statusCode, err := c.Invoke(
			ctx,
			lambdaFunctionName,
			map[string]interface{}{"chat_id": chatID, "action": "on_demand"},
		)
		if err != nil {
			sugar.Errorw(
				"error invoking lambda function",
				"function_name",
				lambdaFunctionName,
				"error",
				err.Error(),
			)
			return createTelegramResponse(http.StatusOK, chatID, "Lo siento, algo ha fallado")
		}

		sugar.Debugw(
			"successfully invoked Lambda function",
			"function_name",
			lambdaFunctionName,
			"chat_id",
			chatID,
			"status_code",
			statusCode,
		)

		return Response{
			Body:       "success",
			StatusCode: http.StatusOK,
		}, nil
	}
	return createTelegramResponse(
		http.StatusOK,
		chatID,
		"Lo siento, solo acepto comandos de Telegram.",
	)
}

func createTelegramResponse(
	status int,
	chatID int64,
	text string,
) (Response, error) {
	t := api.TelegramWebhookSendMessage{
		ChatID: chatID,
		Text:   text,
		Method: "sendMessage",
	}
	headers := map[string]string{
		"Content-Type": "application/json",
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
