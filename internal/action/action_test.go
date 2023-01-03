package action

import (
	"context"
	"lamver/internal/types"
	"lamver/pkg/client"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/golang/mock/gomock"
)

func TestGetAllRegionsAndRuntime(t *testing.T) {
	type args struct {
		ctx    context.Context
		region string
	}
	ctx := context.Background()

	tests := []struct {
		name                          string
		args                          args
		prepareMockAWSConfigCreatorFn func(m *client.MockAWSConfigCreator)
		prepareMockEC2CreatorFn       func(m *client.MockEC2Creator, c *client.MockEC2Client)
		prepareMockLambdaCreatorFn    func(m *client.MockLambdaCreator, c *client.MockLambdaClient)
		wantRegionList                []string
		wantRuntimeList               []string
		wantErr                       bool
	}{
		{
			name: "GetAllRegionsAndRuntime success",
			args: args{
				ctx:    ctx,
				region: "us-east-1",
			},
			prepareMockAWSConfigCreatorFn: func(m *client.MockAWSConfigCreator) {
				m.EXPECT().Create(ctx, "us-east-1").Return(
					&client.AWSConfig{
						Config: aws.Config{},
					}, nil,
				)
			},
			prepareMockEC2CreatorFn: func(m *client.MockEC2Creator, c *client.MockEC2Client) {
				c.EXPECT().DescribeRegions(ctx).Return(
					[]string{
						"ap-northeast-1",
						"us-east-1",
					}, nil,
				)
				m.EXPECT().Create(aws.Config{}).Return(
					c,
				)
			},
			prepareMockLambdaCreatorFn: func(m *client.MockLambdaCreator, c *client.MockLambdaClient) {
				c.EXPECT().ListRuntimeValues().Return(
					[]string{
						"go1.x",
						"nodejs18.x",
					},
				)
				m.EXPECT().Create(aws.Config{}).Return(
					c,
				)
			},
			wantRegionList: []string{
				"ap-northeast-1",
				"us-east-1",
			},
			wantRuntimeList: []string{
				"go1.x",
				"nodejs18.x",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			awsConfigMock := client.NewMockAWSConfigCreator(ctrl)
			ec2FactoryMock := client.NewMockEC2Creator(ctrl)
			lambdaFactoryMock := client.NewMockLambdaCreator(ctrl)
			ec2ClientMock := client.NewMockEC2Client(ctrl)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockAWSConfigCreatorFn(awsConfigMock)
			tt.prepareMockEC2CreatorFn(ec2FactoryMock, ec2ClientMock)
			tt.prepareMockLambdaCreatorFn(lambdaFactoryMock, lambdaClientMock)

			input := &GetAllRegionsAndRuntimeInput{
				Ctx:              tt.args.ctx,
				AWSConfigFactory: awsConfigMock,
				EC2Factory:       ec2FactoryMock,
				LambdaFactory:    lambdaFactoryMock,
				DefaultRegion:    tt.args.region,
			}

			gotRegionList, gotRuntimeList, err := GetAllRegionsAndRuntime(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllRegionsAndRuntime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotRegionList, tt.wantRegionList) {
				t.Errorf("GetAllRegionsAndRuntime() gotRegionList = %v, want %v", gotRegionList, tt.wantRegionList)
			}
			if !reflect.DeepEqual(gotRuntimeList, tt.wantRuntimeList) {
				t.Errorf("GetAllRegionsAndRuntime() gotRuntimeList = %v, want %v", gotRuntimeList, tt.wantRuntimeList)
			}
		})
	}
}

func TestCreateFunctionMap(t *testing.T) {
	type args struct {
		input *CreateFunctionMapInput
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string][][]string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateFunctionMap(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFunctionMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateFunctionMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_putToFunctionChannelByRegion(t *testing.T) {
	type args struct {
		ctx              context.Context
		region           string
		targetRuntime    []string
		keyword          string
		functionCh       chan *types.LambdaFunctionData
		awsConfigFactory client.AWSConfigCreator
		lambdaFactory    client.LambdaCreator
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := putToFunctionChannelByRegion(tt.args.ctx, tt.args.region, tt.args.targetRuntime, tt.args.keyword, tt.args.functionCh, tt.args.awsConfigFactory, tt.args.lambdaFactory); (err != nil) != tt.wantErr {
				t.Errorf("putToFunctionChannelByRegion() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSortAndSetFunctionList(t *testing.T) {
	type args struct {
		input *SortAndSetFunctionListInput
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SortAndSetFunctionList(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortAndSetFunctionList() = %v, want %v", got, tt.want)
			}
		})
	}
}
