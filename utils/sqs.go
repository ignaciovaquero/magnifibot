package utils

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

func InitSQSClient(region, url string) (*sqs.Client, error) {
	var cfg aws.Config
	var err error

	if url == "" {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
		)
	} else {
		customResolver := aws.EndpointResolverWithOptionsFunc(func(service, awsRegion string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           url,
				SigningRegion: awsRegion,
			}, nil
		})
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion(region),
			config.WithEndpointResolverWithOptions(customResolver),
		)
	}
	if err != nil {
		return nil, fmt.Errorf("error loading aws configuration: %w", err)
	}
	return sqs.NewFromConfig(cfg), nil
}
