//go:generate mockgen -source=./lambda.go -destination=./lambda_mock.go -package=client
package client

import (
	"context"
	"sort"

	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

type LambdaClient interface {
	ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error)
	ListFunctionsWithRegion(ctx context.Context, region string) ([]types.FunctionConfiguration, error)
	ListRuntimeValues() []string
}

type LambdaSDKClient interface {
	ListFunctions(ctx context.Context, params *lambda.ListFunctionsInput, optFns ...func(*lambda.Options)) (*lambda.ListFunctionsOutput, error)
}

type Lambda struct {
	client LambdaSDKClient
}

var _ LambdaClient = (*Lambda)(nil)

func NewLambda(client LambdaSDKClient) *Lambda {
	return &Lambda{
		client: client,
	}
}

func (c *Lambda) ListFunctions(ctx context.Context) ([]types.FunctionConfiguration, error) {
	return c.ListFunctionsWithRegion(ctx, "")
}

func (c *Lambda) ListFunctionsWithRegion(ctx context.Context, region string) ([]types.FunctionConfiguration, error) {
	var nextMarker *string
	outputs := []types.FunctionConfiguration{}

	var optFns func(*lambda.Options)
	if region != "" {
		optFns = func(o *lambda.Options) {
			o.Region = region
		}
	}

	for {
		input := &lambda.ListFunctionsInput{
			Marker: nextMarker,
		}

		var (
			output *lambda.ListFunctionsOutput
			err    error
		)

		if region == "" {
			output, err = c.client.ListFunctions(ctx, input)
		} else {
			output, err = c.client.ListFunctions(ctx, input, optFns)
		}
		if err != nil {
			return outputs, err
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
