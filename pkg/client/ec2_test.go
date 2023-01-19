package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go/middleware"
)

func TestEC2_DescribeRegions(t *testing.T) {
	type args struct {
		ctx                context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name: "DescribeRegions success",
			args: args{
				ctx: context.Background(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DescribeRegionsMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &ec2.DescribeRegionsOutput{
										Regions: []types.Region{
											{
												RegionName: aws.String("ap-northeast-1"),
											},
											{
												RegionName: aws.String("us-east-1"),
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
			want: []string{
				"ap-northeast-1",
				"us-east-1",
			},
			wantErr: false,
		},
		{
			name: "DescribeRegions sorted success",
			args: args{
				ctx: context.Background(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DescribeRegionsUnSortedMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &ec2.DescribeRegionsOutput{
										Regions: []types.Region{
											{
												RegionName: aws.String("us-east-1"),
											},
											{
												RegionName: aws.String("ap-northeast-1"),
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
			want: []string{ // sort by region name
				"ap-northeast-1",
				"us-east-1",
			},
			wantErr: false,
		},
		{
			name: "DescribeRegions fail",
			args: args{
				ctx: context.Background(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DescribeRegionsErrorMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &ec2.DescribeRegionsOutput{},
								}, middleware.Metadata{}, fmt.Errorf("DescribeRegionsError")
							},
						),
						middleware.Before,
					)
				},
			},
			want:    []string{},
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

			client := ec2.NewFromConfig(cfg)
			ec2Client := NewEC2(client)

			got, err := ec2Client.DescribeRegions(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("EC2.DescribeRegions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2.DescribeRegions() = %v, want %v", got, tt.want)
			}
		})
	}
}
