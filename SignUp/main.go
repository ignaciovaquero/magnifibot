package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/igvaquero18/magnifibot/api"
	"github.com/igvaquero18/magnifibot/controller"
	"github.com/igvaquero18/magnifibot/utils"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	awsRegionEnv         = "MAGNIFIBOT_AWS_REGION"
	verboseEnv           = "MAGNIFIBOT_VERBOSE"
	dynamoDBEndpointEnv  = "MAGNIFIBOT_DYNAMODB_ENDPOINT"
	dynamoDBUserTableEnv = "MAGNIFIBOT_DYNAMODB_USER_TABLE"
)

const (
	awsRegionFlag         = "aws.region"
	verboseFlag           = "logging.verbose"
	dynamoDBEndpointFlag  = "aws.dynamodb.endpoint"
	dynamoDBUserTableFlag = "aws.dynamodb.tables.user"
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
	viper.SetDefault(awsRegionFlag, "us-east-3")
	viper.SetDefault(verboseFlag, false)
	viper.SetDefault(dynamoDBEndpointFlag, "")
	viper.SetDefault(dynamoDBUserTableFlag, controller.DefaultUserTable)
	viper.BindEnv(awsRegionFlag, awsRegionEnv)
	viper.BindEnv(verboseFlag, verboseEnv)
	viper.BindEnv(dynamoDBEndpointFlag, dynamoDBEndpointEnv)
	viper.BindEnv(dynamoDBUserTableFlag, dynamoDBUserTableEnv)

	sugar, err := utils.InitSugaredLogger(viper.GetBool(verboseFlag))

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

	c = controller.NewMagnifibot(controller.SetDynamoDBClient(dynamoClient), controller.SetConfig(&controller.MagnifibotConfig{
		UserTable: viper.GetString(dynamoDBUserTableFlag),
	}))
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
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
		re := regexp.MustCompile(`/(\w*)@?\w*`)
		if !re.Match([]byte(update.Message.Text)) {
			body, err := json.Marshal(map[string]interface{}{
				"text": "Invalid command",
			})
			if err != nil {
				return Response{
					Body:       fmt.Sprintf("error when marshalling response: %s", err.Error()),
					StatusCode: http.StatusInternalServerError,
					Headers:    headers,
				}, fmt.Errorf("error when marshalling response: %w", err)
			}
			return Response{
				Body:       string(body),
				StatusCode: http.StatusBadRequest,
				Headers:    headers,
			}, nil
		}
		submatch := re.FindStringSubmatch(update.Message.Text)
		text := ""
		if len(submatch) > 1 {
			text = submatch[1]
		}
	}
}

func main() {
	lambda.Start(Handler)
}
