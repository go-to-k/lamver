package action

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"

	"github.com/go-to-k/lamver/internal/types"
	"github.com/go-to-k/lamver/pkg/client"

	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"go.uber.org/mock/gomock"
)

func TestGetAllRegionsAndRuntime(t *testing.T) {
	type args struct {
		ctx    context.Context
		region string
	}

	tests := []struct {
		name                      string
		args                      args
		prepareMockEC2ClientFn    func(m *client.MockEC2Client)
		prepareMockLambdaClientFn func(m *client.MockLambdaClient)
		wantRegionList            []string
		wantRuntimeList           []string
		wantErr                   bool
	}{
		{
			name: "GetAllRegionsAndRuntime success",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
			},
			prepareMockEC2ClientFn: func(m *client.MockEC2Client) {
				m.EXPECT().DescribeRegions(gomock.Any()).Return(
					[]string{
						"ap-northeast-1",
						"us-east-1",
					}, nil,
				)
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListRuntimeValues().Return(
					[]string{
						"go1.x",
						"nodejs18.x",
					},
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
		{
			name: "GetAllRegionsAndRuntime fail by DescribeRegions Error",
			args: args{
				ctx:    context.Background(),
				region: "us-east-1",
			},
			prepareMockEC2ClientFn: func(m *client.MockEC2Client) {
				m.EXPECT().DescribeRegions(gomock.Any()).Return(
					[]string{}, fmt.Errorf("DescribeRegionsError"),
				)
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListRuntimeValues().Return(
					[]string{
						"go1.x",
						"nodejs18.x",
					},
				)
			},
			wantRegionList: []string{},
			wantRuntimeList: []string{
				"go1.x",
				"nodejs18.x",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ec2ClientMock := client.NewMockEC2Client(ctrl)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockEC2ClientFn(ec2ClientMock)
			tt.prepareMockLambdaClientFn(lambdaClientMock)

			input := &GetAllRegionsAndRuntimeInput{
				Ctx:           tt.args.ctx,
				EC2:           ec2ClientMock,
				Lambda:        lambdaClientMock,
				DefaultRegion: tt.args.region,
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

func TestCreateFunctionList(t *testing.T) {
	type args struct {
		ctx           context.Context
		targetRegions []string
		targetRuntime []string
		keyword       string
	}

	tests := []struct {
		name                      string
		args                      args
		prepareMockLambdaClientFn func(m *client.MockLambdaClient)
		want                      [][]string
		wantErr                   bool
	}{
		{
			name: "CreateFunctionList success",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs18.x", "us-east-2", "Function6", "2022-12-22T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "Function3", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-2", "Function5", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success but there is no function",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{}, nil,
				)
			},
			want:    [][]string{},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if a keyword is not empty",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "3",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs", "us-east-1", "Function3", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if a keyword is not empty and lower",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "function3",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs", "us-east-1", "Function3", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if a keyword is not empty and upper",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "FUNCTION3",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs", "us-east-1", "Function3", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if a keyword is not empty when except runtime",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "2",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want:    [][]string{},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if any regions have empty function list",
			args: args{
				ctx:           context.Background(),
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs18.x", "us-east-2", "Function6", "2022-12-22T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "Function3", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-2", "Function5", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockLambdaClientFn(lambdaClientMock)

			input := &CreateFunctionListInput{
				Ctx:           tt.args.ctx,
				TargetRegions: tt.args.targetRegions,
				TargetRuntime: tt.args.targetRuntime,
				Keyword:       tt.args.keyword,
				Lambda:        lambdaClientMock,
			}

			got, err := CreateFunctionList(input)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateFunctionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && (len(got) != 0 || len(tt.want) != 0) {
				t.Errorf("CreateFunctionList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_putToFunctionChannelByRegion(t *testing.T) {
	type args struct {
		ctx           context.Context
		region        string
		targetRuntime []string
		keyword       string
		functionCh    chan *types.LambdaFunctionData
	}

	tests := []struct {
		name                      string
		args                      args
		prepareMockLambdaClientFn func(m *client.MockLambdaClient)
		putCount                  int
		wantErr                   bool
	}{
		{
			name: "putToFunctionChannelByRegion success",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 1,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success if there is no corresponding runtime",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x"},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 0,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success if lower keywords given",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "function3",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 1,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success if upper keywords given",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "FUNCTION3",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 1,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success if there is no corresponding function matching the given region keywords",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "any",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 0,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success but there is no function",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{}, nil,
				)
			},
			putCount: 0,
			wantErr:  false,
		},
		{
			name: "putToFunctionChannelByRegion success but there is no targetRuntime",
			args: args{
				ctx:           context.Background(),
				region:        "us-east-1",
				targetRuntime: []string{},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("Function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("Function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			putCount: 0,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockLambdaClientFn(lambdaClientMock)

			putCount := 0
			ctx, cancel := context.WithCancel(tt.args.ctx)
			ch := tt.args.functionCh
			wg := sync.WaitGroup{}

			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ch:
						putCount++
					}
				}
			}()

			if err := putToFunctionChannelByRegion(ctx, tt.args.region, tt.args.targetRuntime, tt.args.keyword, ch, lambdaClientMock); (err != nil) != tt.wantErr {
				t.Errorf("putToFunctionChannelByRegion() error = %v, wantErr %v", err, tt.wantErr)
				cancel()
				return
			}
			cancel()
			wg.Wait()
			if !tt.wantErr && putCount != tt.putCount {
				t.Errorf("putCount = %v, tt.putCount %v", putCount, tt.putCount)
			}
		})
	}
}

func Test_sortAndSetFunctionList(t *testing.T) {
	type args struct {
		regionList  []string
		runtimeList []string
		functionMap map[string]map[string][][]string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "sortAndSetFunctionList success",
			args: args{
				regionList:  []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				runtimeList: []string{"nodejs", "nodejs18.x"},
				functionMap: map[string]map[string][][]string{
					"nodejs": {
						"ap-northeast-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
						"us-east-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
					},
					"nodejs18.x": {
						"ap-northeast-1": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
						"us-east-2": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
					},
				},
			},
			want: [][]string{
				{"nodejs", "ap-northeast-1", "Function1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "Function1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs18.x", "ap-northeast-1", "Function3", "2022-12-22T09:47:43.728+0000"},
				{"nodejs18.x", "us-east-2", "Function3", "2022-12-22T09:47:43.728+0000"},
			},
		},
		{
			name: "sortAndSetFunctionList success with sorted by function names per regions and runtime",
			args: args{
				regionList:  []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				runtimeList: []string{"nodejs", "nodejs18.x"},
				functionMap: map[string]map[string][][]string{
					"nodejs": {
						"ap-northeast-1": {
							[]string{"Function-A", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-c", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-b", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-B", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-a", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-1", "2022-12-21T09:47:43.728+0000"},
						},
						"us-east-1": {
							[]string{"Function-b-1", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-a-2", "2022-12-21T09:47:43.728+0000"},
						},
					},
					"nodejs18.x": {
						"ap-northeast-1": {
							[]string{"Function-a-3", "2022-12-21T09:47:43.728+0000"},
							[]string{"Function-a-0", "2022-12-21T09:47:43.728+0000"},
						},
					},
				},
			},
			want: [][]string{
				{"nodejs", "ap-northeast-1", "Function", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-A", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-B", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-a", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-b", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "Function-c", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "Function-a-2", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "Function-b-1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs18.x", "ap-northeast-1", "Function-a-0", "2022-12-21T09:47:43.728+0000"},
				{"nodejs18.x", "ap-northeast-1", "Function-a-3", "2022-12-21T09:47:43.728+0000"},
			},
		},
		{
			name: "sortAndSetFunctionList success if runtimeList is empty",
			args: args{
				regionList:  []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				runtimeList: []string{},
				functionMap: map[string]map[string][][]string{
					"nodejs": {
						"ap-northeast-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
						"us-east-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
					},
					"nodejs18.x": {
						"ap-northeast-1": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
						"us-east-2": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
					},
				},
			},
			want: [][]string{},
		},
		{
			name: "sortAndSetFunctionList success if regionList is empty",
			args: args{
				regionList:  []string{},
				runtimeList: []string{"nodejs", "nodejs18.x"},
				functionMap: map[string]map[string][][]string{
					"nodejs": {
						"ap-northeast-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
						"us-east-1": {
							[]string{"Function1", "2022-12-21T09:47:43.728+0000"},
						},
					},
					"nodejs18.x": {
						"ap-northeast-1": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
						"us-east-2": {
							[]string{"Function3", "2022-12-22T09:47:43.728+0000"},
						},
					},
				},
			},
			want: [][]string{},
		},
		{
			name: "sortAndSetFunctionList success if functionMap is empty",
			args: args{
				regionList:  []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				runtimeList: []string{"nodejs", "nodejs18.x"},
				functionMap: map[string]map[string][][]string{},
			},
			want: [][]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortAndSetFunctionList(tt.args.regionList, tt.args.runtimeList, tt.args.functionMap)
			if !reflect.DeepEqual(got, tt.want) && (len(got) != 0 || len(tt.want) != 0) {
				t.Errorf("sortAndSetFunctionList() = %v, want %v", got, tt.want)
			}
		})
	}
}
