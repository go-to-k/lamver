package client

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
)

func TestNewLambdaClient(t *testing.T) {
	type args struct {
		config aws.Config
	}
	tests := []struct {
		name string
		args args
		want *Lambda
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLambdaClient(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLambdaClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLambda_ListFunctions(t *testing.T) {
	type fields struct {
		client *lambda.Client
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []types.FunctionConfiguration
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Lambda{
				client: tt.fields.client,
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
		client *lambda.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
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
