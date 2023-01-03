package client

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/golang/mock/gomock"
)

func TestNewEC2(t *testing.T) {
	type args struct {
		client EC2SDKClient
	}
	ctrl := gomock.NewController(t)
	mock := NewMockEC2SDKClient(ctrl)

	tests := []struct {
		name string
		args args
		want *EC2
	}{
		{
			name: "NewEC2",
			args: args{
				client: mock,
			},
			want: &EC2{
				client: mock,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEC2(tt.args.client); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEC2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEC2_DescribeRegions(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	ctx := context.Background()

	tests := []struct {
		name          string
		args          args
		prepareMockFn func(m *MockEC2SDKClient)
		want          []string
		wantErr       bool
	}{
		{
			name: "DescribeRegions success",
			args: args{
				ctx: ctx,
			},
			prepareMockFn: func(m *MockEC2SDKClient) {
				m.EXPECT().DescribeRegions(ctx, &ec2.DescribeRegionsInput{}).Return(
					&ec2.DescribeRegionsOutput{
						Regions: []types.Region{
							{
								RegionName: aws.String("ap-northeast-1"),
							},
							{
								RegionName: aws.String("us-east-1"),
							},
						},
					}, nil,
				)
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
				ctx: ctx,
			},
			prepareMockFn: func(m *MockEC2SDKClient) {
				m.EXPECT().DescribeRegions(ctx, &ec2.DescribeRegionsInput{}).Return(
					&ec2.DescribeRegionsOutput{
						Regions: []types.Region{
							{
								RegionName: aws.String("us-east-1"),
							},
							{
								RegionName: aws.String("ap-northeast-1"),
							}},
					}, nil,
				)
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
				ctx: ctx,
			},
			prepareMockFn: func(m *MockEC2SDKClient) {
				m.EXPECT().DescribeRegions(ctx, &ec2.DescribeRegionsInput{}).Return(
					nil, fmt.Errorf("DescribeRegionsError"),
				)
			},
			want:    []string{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := NewMockEC2SDKClient(ctrl)

			tt.prepareMockFn(mock)

			c := &EC2{
				client: mock,
			}
			got, err := c.DescribeRegions(tt.args.ctx)
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
