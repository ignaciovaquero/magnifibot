package utils

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/stretchr/testify/assert"
)

func TestInitDynamoClient(t *testing.T) {
	cfg, _ := config.LoadDefaultConfig(context.TODO())
	testCases := []struct {
		name          string
		region        string
		url           string
		expectedError bool
	}{
		{
			name:          "Empty URL, eu-west-3",
			region:        "eu-west-3",
			url:           "",
			expectedError: false,
		},
		{
			name:          "Fake region",
			region:        "fake",
			url:           "",
			expectedError: false,
		},
		{
			name:          "Localhost url, eu-west-3",
			region:        "eu-west-3",
			url:           "http://localhost:8000",
			expectedError: false,
		},
		{
			name:          "Invalid url",
			region:        "eu-west-3",
			url:           "invalid url",
			expectedError: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(tt *testing.T) {
			cfg.Region = tc.region
			if tc.url != "" {
				cfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(func(service, awsRegion string, options ...interface{}) (aws.Endpoint, error) {
					return aws.Endpoint{
						PartitionID:   "aws",
						URL:           tc.url,
						SigningRegion: awsRegion,
					}, nil
				})
			}
			expected := dynamodb.NewFromConfig(cfg)
			actual, err := InitDynamoClient(tc.region, tc.url)
			if tc.expectedError {
				assert.Error(tt, err)
				return
			}
			assert.NoError(tt, err)
			assert.ObjectsAreEqual(expected, actual)
		})
	}
}
