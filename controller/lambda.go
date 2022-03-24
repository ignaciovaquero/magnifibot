package controller

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func (m *Magnifibot) Invoke(ctx context.Context, functionName string, payload map[string]interface{}) (int32, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return 0, fmt.Errorf("error converting payload to JSON: %w", err)
	}
	output, err := m.LambdaInvokeAPI.Invoke(ctx, &lambda.InvokeInput{
		FunctionName:   aws.String(functionName),
		InvocationType: types.InvocationTypeEvent,
		Payload:        body,
	})
	if err != nil {
		return 0, fmt.Errorf("error invoking lambda function %q: %w", functionName, err)
	}
	return output.StatusCode, nil
}
