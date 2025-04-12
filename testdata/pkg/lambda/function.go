package lambda

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

const (
	FunctionName = "lamver-test-function"
)

// Helper function to return an int32 pointer
func Int32(v int32) *int32 {
	return &v
}

// Create Lambda function
func CreateFunction(ctx context.Context, cfg aws.Config, funcName, roleARN string, runtimeInfo RuntimeInfo) (bool, error) {
	// Create Lambda client
	client := lambda.NewFromConfig(cfg)

	// Check if function already exists
	_, err := client.GetFunction(ctx, &lambda.GetFunctionInput{
		FunctionName: aws.String(funcName),
	})

	// If function exists, return without error
	if err == nil {
		// Return false to indicate the function already existed
		return false, nil
	}

	// Compress Lambda function source code into a ZIP file
	zipBytes, err := createZip(runtimeInfo)
	if err != nil {
		return false, fmt.Errorf("failed to create ZIP file: %w", err)
	}

	// Create Lambda function
	_, err = client.CreateFunction(ctx, &lambda.CreateFunctionInput{
		Code: &types.FunctionCode{
			ZipFile: zipBytes,
		},
		FunctionName: aws.String(funcName),
		Handler:      aws.String(runtimeInfo.Handler),
		Role:         aws.String(roleARN),
		Runtime:      runtimeInfo.Runtime,
		Timeout:      Int32(30),
		MemorySize:   Int32(128),
	})

	if err != nil {
		return false, err
	}

	// Return true to indicate the function was newly created
	return true, nil
}

// Delete Lambda function
func DeleteFunction(ctx context.Context, cfg aws.Config, funcName string) error {
	// Create Lambda client
	client := lambda.NewFromConfig(cfg)

	// Delete Lambda function
	_, err := client.DeleteFunction(ctx, &lambda.DeleteFunctionInput{
		FunctionName: aws.String(funcName),
	})

	return err
}
