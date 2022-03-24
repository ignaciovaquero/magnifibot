package controller

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/stretchr/testify/assert"
)

type MockLambda struct {
	statusCode int32
	err        error
}

func (m *MockLambda) Invoke(
	ctx context.Context,
	params *lambda.InvokeInput,
	optFns ...func(*lambda.Options),
) (*lambda.InvokeOutput, error) {
	return &lambda.InvokeOutput{StatusCode: m.statusCode}, m.err
}

func TestInvoke(t *testing.T) {
	tests := []struct {
		name,
		functionName string
		payload       map[string]interface{}
		expected      int32
		client        LambdaInvokeAPI
		errorExpected bool
	}{
		{
			name:         "valid invocation",
			functionName: "lambda-function",
			payload:      map[string]interface{}{"test": "test"},
			expected:     http.StatusAccepted,
			client: &MockLambda{
				statusCode: http.StatusAccepted,
				err:        nil,
			},
			errorExpected: false,
		},
		{
			name:         "non-existing function",
			functionName: "lambda-non-existing",
			payload:      map[string]interface{}{"test": "test"},
			expected:     0,
			client: &MockLambda{
				statusCode: http.StatusBadRequest,
				err:        errors.New("error"),
			},
			errorExpected: true,
		},
		{
			name:         "valid invocation with integer value",
			functionName: "lambda-function",
			payload:      map[string]interface{}{"test": 0},
			expected:     http.StatusAccepted,
			client: &MockLambda{
				statusCode: http.StatusAccepted,
				err:        nil,
			},
			errorExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			m := NewMagnifibot(SetLambdaClient(test.client))
			actual, err := m.Invoke(context.TODO(), test.functionName, test.payload)
			if test.errorExpected {
				assert.Error(tt, err)
			} else {
				assert.NoError(tt, err)
			}
			assert.Equal(tt, test.expected, actual)
		})
	}
}
