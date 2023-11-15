package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/aws/smithy-go/middleware"
)

type markerKey struct{}

func getNextMarkerForInitialize(
	ctx context.Context, in middleware.InitializeInput, next middleware.InitializeHandler,
) (
	out middleware.InitializeOutput, metadata middleware.Metadata, err error,
) {
	switch v := in.Parameters.(type) {
	case *lambda.ListFunctionsInput:
		ctx = middleware.WithStackValue(ctx, markerKey{}, v.Marker)
	}
	return next.HandleInitialize(ctx, in)
}

func TestLambda_ListFunctionsWithRegion(t *testing.T) {
	type args struct {
		ctx                context.Context
		region             string
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		args    args
		want    []types.FunctionConfiguration
		wantErr bool
	}{
		{
			name: "ListFunctionsWithRegion success",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{
										NextMarker: nil,
										Functions: []types.FunctionConfiguration{
											{
												FunctionName: aws.String("Function1"),
												Runtime:      types.RuntimeNodejs,
												LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
											},
											{
												FunctionName: aws.String("Function2"),
												Runtime:      types.RuntimeNodejs18x,
												LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
											},
										},
									},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion with no functions success",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNoFunctionsMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{
										NextMarker: nil,
										Functions:  []types.FunctionConfiguration{},
									},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			want:    []types.FunctionConfiguration{},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion with NextMarker success",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					err := stack.Initialize.Add(
						middleware.InitializeMiddlewareFunc(
							"GetNextMarkerFromListFunctionsInput",
							getNextMarkerForInitialize,
						), middleware.Before,
					)
					if err != nil {
						return err
					}

					err = stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNextMarkerMock",
							func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								marker := middleware.GetStackValue(ctx, markerKey{}).(*string)

								var nextMarker *string
								var functions []types.FunctionConfiguration
								if marker == nil {
									nextMarker = aws.String("NextMarker")
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function1"),
											Runtime:      types.RuntimeNodejs,
											LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function2"),
											Runtime:      types.RuntimeNodejs18x,
											LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								} else {
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function3"),
											Runtime:      types.RuntimeGo1x,
											LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function4"),
											Runtime:      types.RuntimeProvidedal2,
											LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								}
							},
						),
						middleware.Before,
					)
					return err
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function3"),
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function4"),
					Runtime:      types.RuntimeProvidedal2,
					LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion fail",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsErrorMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{},
								}, middleware.Metadata{}, fmt.Errorf("ListFunctionsError")
							},
						),
						middleware.Before,
					)
				},
			},
			want:    []types.FunctionConfiguration{},
			wantErr: true,
		},
		{
			name: "ListFunctionsWithRegion with NextMarker fail",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					err := stack.Initialize.Add(
						middleware.InitializeMiddlewareFunc(
							"GetNextMarkerFromListFunctionsInput",
							getNextMarkerForInitialize,
						), middleware.Before,
					)
					if err != nil {
						return err
					}

					err = stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNextMarkerErrorMock",
							func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								marker := middleware.GetStackValue(ctx, markerKey{}).(*string)

								var nextMarker *string
								var functions []types.FunctionConfiguration
								if marker == nil {
									nextMarker = aws.String("NextMarker")
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function1"),
											Runtime:      types.RuntimeNodejs,
											LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function2"),
											Runtime:      types.RuntimeNodejs18x,
											LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								} else {
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{},
									}, middleware.Metadata{}, fmt.Errorf("ListFunctionsError")
								}
							},
						),
						middleware.Before,
					)
					return err
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: true,
		},
		{
			name: "ListFunctionsWithRegion with empty region success",
			args: args{
				ctx:    context.Background(),
				region: "",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithEmptyRegionMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{
										NextMarker: nil,
										Functions: []types.FunctionConfiguration{
											{
												FunctionName: aws.String("Function1"),
												Runtime:      types.RuntimeNodejs,
												LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
											},
											{
												FunctionName: aws.String("Function2"),
												Runtime:      types.RuntimeNodejs18x,
												LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
											},
										},
									},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion with no functions and empty region success",
			args: args{
				ctx:    context.Background(),
				region: "",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNoFunctionsAndEmptyRegionMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{
										NextMarker: nil,
										Functions:  []types.FunctionConfiguration{},
									},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			want:    []types.FunctionConfiguration{},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion with NextMarker and empty region success",
			args: args{
				ctx:    context.Background(),
				region: "",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					err := stack.Initialize.Add(
						middleware.InitializeMiddlewareFunc(
							"GetNextMarkerFromListFunctionsInput",
							getNextMarkerForInitialize,
						), middleware.Before,
					)
					if err != nil {
						return err
					}

					err = stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNextMarkerAndEmptyRegionMock",
							func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								marker := middleware.GetStackValue(ctx, markerKey{}).(*string)

								var nextMarker *string
								var functions []types.FunctionConfiguration
								if marker == nil {
									nextMarker = aws.String("NextMarker")
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function1"),
											Runtime:      types.RuntimeNodejs,
											LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function2"),
											Runtime:      types.RuntimeNodejs18x,
											LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								} else {
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function3"),
											Runtime:      types.RuntimeGo1x,
											LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function4"),
											Runtime:      types.RuntimeProvidedal2,
											LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								}
							},
						),
						middleware.Before,
					)
					return err
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function3"),
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function4"),
					Runtime:      types.RuntimeProvidedal2,
					LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctionsWithRegion with empty region fail",
			args: args{
				ctx:    context.Background(),
				region: "",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithEmptyRegionErrorMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &lambda.ListFunctionsOutput{
										NextMarker: nil,
										Functions:  []types.FunctionConfiguration{},
									},
								}, middleware.Metadata{}, fmt.Errorf("ListFunctionsError")
							},
						),
						middleware.Before,
					)
				},
			},
			want:    []types.FunctionConfiguration{},
			wantErr: true,
		},
		{
			name: "ListFunctionsWithRegion with NextMarker and empty region fail",
			args: args{
				ctx:    context.Background(),
				region: "",
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					err := stack.Initialize.Add(
						middleware.InitializeMiddlewareFunc(
							"GetNextMarkerFromListFunctionsInput",
							getNextMarkerForInitialize,
						), middleware.Before,
					)
					if err != nil {
						return err
					}

					err = stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"ListFunctionsWithNextMarkerAndEmptyRegionErrorMock",
							func(ctx context.Context, input middleware.FinalizeInput, handler middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								marker := middleware.GetStackValue(ctx, markerKey{}).(*string)

								var nextMarker *string
								var functions []types.FunctionConfiguration
								if marker == nil {
									nextMarker = aws.String("NextMarker")
									functions = []types.FunctionConfiguration{
										{
											FunctionName: aws.String("Function1"),
											Runtime:      types.RuntimeNodejs,
											LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
										},
										{
											FunctionName: aws.String("Function2"),
											Runtime:      types.RuntimeNodejs18x,
											LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
										},
									}
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{
											NextMarker: nextMarker,
											Functions:  functions,
										},
									}, middleware.Metadata{}, nil
								} else {
									return middleware.FinalizeOutput{
										Result: &lambda.ListFunctionsOutput{},
									}, middleware.Metadata{}, fmt.Errorf("ListFunctionsError")
								}
							},
						),
						middleware.Before,
					)
					return err
				},
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("Function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("Function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(
				tt.args.ctx,
				config.WithRegion("ap-northeast-1"),
				config.WithAPIOptions([]func(*middleware.Stack) error{tt.args.withAPIOptionsFunc}),
			)
			if err != nil {
				t.Fatal(err)
			}

			client := lambda.NewFromConfig(cfg)
			lambdaClient := NewLambda(client)

			got, err := lambdaClient.ListFunctionsWithRegion(tt.args.ctx, tt.args.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lambda.ListFunctionsWithRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lambda.ListFunctionsWithRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambda_ListRuntimeValues(t *testing.T) {
	tests := []struct {
		name string
		want []string
	}{
		{
			name: "ListRuntimeValues sorted success",
			want: []string{
				"dotnet6",
				"dotnetcore1.0",
				"dotnetcore2.0",
				"dotnetcore2.1",
				"dotnetcore3.1",
				"go1.x",
				"java8",
				"java8.al2",
				"java11",
				"java17",
				"nodejs",
				"nodejs4.3",
				"nodejs4.3-edge",
				"nodejs6.10",
				"nodejs8.10",
				"nodejs10.x",
				"nodejs12.x",
				"nodejs14.x",
				"nodejs16.x",
				"nodejs18.x",
				"nodejs20.x",
				"provided",
				"provided.al2",
				"provided.al2023",
				"python2.7",
				"python3.6",
				"python3.7",
				"python3.8",
				"python3.9",
				"python3.10",
				"python3.11",
				"python3.12",
				"ruby2.5",
				"ruby2.7",
				"ruby3.2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(
				context.Background(),
				config.WithRegion("ap-northeast-1"),
			)
			if err != nil {
				t.Fatal(err)
			}

			client := lambda.NewFromConfig(cfg)
			lambdaClient := NewLambda(client)

			if got := lambdaClient.ListRuntimeValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lambda.ListRuntimeValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
