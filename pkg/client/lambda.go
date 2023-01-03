package client

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaClient interface {
	ListFunctions(ctx context.Context) ([]string, error)
	ListRuntimeValues() []string
}

type LambdaSDKClient interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

type Lambda struct {
	client LambdaSDKClient
}

func NewLambda(client LambdaSDKClient) *Lambda {
	return &Lambda{
		client: client,
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

func (c *Lambda) ListRuntimeValues() []string {
	var r types.Runtime
	runtimeStrList := []string{}
	runtimeList := r.Values()

	for _, runtime := range runtimeList {
		runtimeStrList = append(runtimeStrList, string(runtime))
	}

	sort.Strings(runtimeStrList)
	return runtimeStrList
}
