package client

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type Lambda struct {
	client *lambda.Client
}

func NewLambdaClient(config aws.Config) *Lambda {
	lambdaClient := lambda.NewFromConfig(config)

	return &Lambda{
		client: lambdaClient,
	}
}

func (c *Lambda) ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error) {
	var nextMarker *string
	outputs := []types.FunctionConfiguration{}

	for {
		input := &lambda.ListFunctionsInput{
			Marker: nextMarker,
		}

		output, err := c.client.ListFunctions(ctx, input)
		if err != nil {
			return nil, err
		}

		outputs = append(outputs, output.Functions...)

		nextMarker = output.NextMarker

		if nextMarker == nil {
			break
		}
	}

	return outputs, nil
}
