package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	gomock "github.com/golang/mock/gomock"
)

func TestNewLambda(t *testing.T) {
	type args struct {
		client LambdaSDKClient
	}
	ctrl := gomock.NewController(t)
	mock := NewMockLambdaSDKClient(ctrl)

	tests := []struct {
		name string
		args args
		want *Lambda
	}{
		{
			name: "NewLambda",
			args: args{
				client: mock,
			},
			want: &Lambda{
				client: mock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLambda(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLambda() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambda_ListFunctions(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name          string
		args          args
		prepareMockFn func(m *MockLambdaSDKClient)
		want          []types.FunctionConfiguration
		wantErr       bool
	}{
		{
			name: "ListFunctions success",
			args: args{
				ctx: ctx,
			},
			prepareMockFn: func(m *MockLambdaSDKClient) {
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: nil}).Return(
					&lambda.ListFunctionsOutput{
						NextMarker: nil,
						Functions: []types.FunctionConfiguration{
							{
								FunctionName: aws.String("function1"),
								Runtime:      types.RuntimeNodejs,
								LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
							},
							{
								FunctionName: aws.String("function2"),
								Runtime:      types.RuntimeNodejs18x,
								LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
							},
						},
					}, nil,
				)
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctions with NextMarker success",
			args: args{
				ctx: ctx,
			},
			prepareMockFn: func(m *MockLambdaSDKClient) {
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: nil}).Return(
					&lambda.ListFunctionsOutput{
						NextMarker: aws.String("NextMarker"),
						Functions: []types.FunctionConfiguration{
							{
								FunctionName: aws.String("function1"),
								Runtime:      types.RuntimeNodejs,
								LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
							},
							{
								FunctionName: aws.String("function2"),
								Runtime:      types.RuntimeNodejs18x,
								LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
							},
						},
					}, nil,
				)
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: aws.String("NextMarker")}).Return(
					&lambda.ListFunctionsOutput{
						NextMarker: nil,
						Functions: []types.FunctionConfiguration{
							{
								FunctionName: aws.String("function3"),
								Runtime:      types.RuntimeGo1x,
								LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
							},
							{
								FunctionName: aws.String("function4"),
								Runtime:      types.RuntimeProvidedal2,
								LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
							},
						},
					}, nil,
				)
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("function3"),
					Runtime:      types.RuntimeGo1x,
					LastModified: aws.String("2022-12-21T10:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("function4"),
					Runtime:      types.RuntimeProvidedal2,
					LastModified: aws.String("2022-12-22T11:47:43.728+0000"),
				},
			},
			wantErr: false,
		},
		{
			name: "ListFunctions fail",
			args: args{
				ctx: ctx,
			},
			prepareMockFn: func(m *MockLambdaSDKClient) {
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: nil}).Return(
					nil, fmt.Errorf("ListFunctionsError"),
				)
			},
			want:    []types.FunctionConfiguration{},
			wantErr: true,
		},
		{
			name: "ListFunctions with NextMarker fail",
			args: args{
				ctx: ctx,
			},
			prepareMockFn: func(m *MockLambdaSDKClient) {
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: nil}).Return(
					&lambda.ListFunctionsOutput{
						NextMarker: aws.String("NextMarker"),
						Functions: []types.FunctionConfiguration{
							{
								FunctionName: aws.String("function1"),
								Runtime:      types.RuntimeNodejs,
								LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
							},
							{
								FunctionName: aws.String("function2"),
								Runtime:      types.RuntimeNodejs18x,
								LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
							},
						},
					}, nil,
				)
				m.EXPECT().ListFunctions(ctx, &lambda.ListFunctionsInput{Marker: aws.String("NextMarker")}).Return(
					nil, fmt.Errorf("ListFunctionsError"),
				)
			},
			want: []types.FunctionConfiguration{
				{
					FunctionName: aws.String("function1"),
					Runtime:      types.RuntimeNodejs,
					LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
				},
				{
					FunctionName: aws.String("function2"),
					Runtime:      types.RuntimeNodejs18x,
					LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := NewMockLambdaSDKClient(ctrl)

			tt.prepareMockFn(mock)

			c := &Lambda{
				client: mock,
			}
			got, err := c.ListFunctions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Lambda.ListFunctions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lambda.ListFunctions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambda_ListRuntimeValues(t *testing.T) {
	type fields struct {
		client LambdaSDKClient
	}
	ctrl := gomock.NewController(t)
	mock := NewMockLambdaSDKClient(ctrl)

	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "ListRuntimeValues sorted success",
			fields: fields{
				client: mock,
			},
			want: []string{
				"dotnet6",
				"dotnetcore1.0",
				"dotnetcore2.0",
				"dotnetcore2.1",
				"dotnetcore3.1",
				"go1.x",
				"java11",
				"java8",
				"java8.al2",
				"nodejs",
				"nodejs10.x",
				"nodejs12.x",
				"nodejs14.x",
				"nodejs16.x",
				"nodejs18.x",
				"nodejs4.3",
				"nodejs4.3-edge",
				"nodejs6.10",
				"nodejs8.10",
				"provided",
				"provided.al2",
				"python2.7",
				"python3.6",
				"python3.7",
				"python3.8",
				"python3.9",
				"ruby2.5",
				"ruby2.7",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Lambda{
				client: tt.fields.client,
			}
			if got := c.ListRuntimeValues(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Lambda.ListRuntimeValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
