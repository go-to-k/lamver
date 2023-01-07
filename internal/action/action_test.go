package action

import (
	"context"
	"fmt"
	"lamver/internal/types"
	"lamver/pkg/client"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	lambdaTypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/golang/mock/gomock"
)

func TestGetAllRegionsAndRuntime(t *testing.T) {
	type args struct {
		ctx    context.Context
		region string
	}
	ctx := context.Background()

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
				ctx:    ctx,
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
				ctx:    ctx,
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
	ctx := context.Background()

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
				ctx:           ctx,
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs18.x", "us-east-2", "function6", "2022-12-22T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "function1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "function3", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-2", "function5", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success but there is no function",
			args: args{
				ctx:           ctx,
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
				ctx:           ctx,
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "3",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs", "us-east-1", "function3", "2022-12-21T09:47:43.728+0000"},
			},
			wantErr: false,
		},
		{
			name: "CreateFunctionList success if a keyword is not empty when except runtime",
			args: args{
				ctx:           ctx,
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "2",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "ap-northeast-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function1"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function2"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function6"),
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
				ctx:           ctx,
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
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-2").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function5"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function6"),
							Runtime:      lambdaTypes.RuntimeNodejs18x,
							LastModified: aws.String("2022-12-22T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			want: [][]string{
				{"nodejs18.x", "us-east-2", "function6", "2022-12-22T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "function3", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-2", "function5", "2022-12-21T09:47:43.728+0000"},
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
			if !reflect.DeepEqual(got, tt.want) {
				if len(got) != len(tt.want) {
					t.Errorf("CreateFunctionList() = %v, want %v", got, tt.want)
				}
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
	ctx := context.Background()

	tests := []struct {
		name                      string
		args                      args
		prepareMockLambdaClientFn func(m *client.MockLambdaClient)
		wantErr                   bool
	}{
		{
			name: "putToFunctionChannelByRegion success",
			args: args{
				ctx:           ctx,
				region:        "us-east-1",
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			wantErr: false,
		},
		{
			name: "putToFunctionChannelByRegion success but there is no function",
			args: args{
				ctx:           ctx,
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
			wantErr: false,
		},
		{
			name: "putToFunctionChannelByRegion success but there is no targetRuntime",
			args: args{
				ctx:           ctx,
				region:        "us-east-1",
				targetRuntime: []string{},
				keyword:       "",
				functionCh:    make(chan *types.LambdaFunctionData),
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(gomock.Any(), "us-east-1").Return(
					[]lambdaTypes.FunctionConfiguration{
						{
							FunctionName: aws.String("function3"),
							Runtime:      lambdaTypes.RuntimeNodejs,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
						{
							FunctionName: aws.String("function4"),
							Runtime:      lambdaTypes.RuntimeGo1x,
							LastModified: aws.String("2022-12-21T09:47:43.728+0000"),
						},
					}, nil,
				)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockLambdaClientFn(lambdaClientMock)

			ch := tt.args.functionCh
			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-ch:
					default:
					}
				}
			}()

			if err := putToFunctionChannelByRegion(ctx, tt.args.region, tt.args.targetRuntime, tt.args.keyword, ch, lambdaClientMock); (err != nil) != tt.wantErr {
				t.Errorf("putToFunctionChannelByRegion() error = %v, wantErr %v", err, tt.wantErr)
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
				runtimeList: []string{"nodejs18.x", "nodejs"},
				functionMap: map[string]map[string][][]string{
					"nodejs": {
						"ap-northeast-1": {
							[]string{"function1", "2022-12-21T09:47:43.728+0000"},
						},
						"us-east-1": {
							[]string{"function1", "2022-12-21T09:47:43.728+0000"},
						},
					},
					"nodejs18.x": {
						"ap-northeast-1": {
							[]string{"function3", "2022-12-22T09:47:43.728+0000"},
						},
						"us-east-2": {
							[]string{"function3", "2022-12-22T09:47:43.728+0000"},
						},
					},
				},
			},
			want: [][]string{
				{"nodejs18.x", "ap-northeast-1", "function3", "2022-12-22T09:47:43.728+0000"},
				{"nodejs18.x", "us-east-2", "function3", "2022-12-22T09:47:43.728+0000"},
				{"nodejs", "ap-northeast-1", "function1", "2022-12-21T09:47:43.728+0000"},
				{"nodejs", "us-east-1", "function1", "2022-12-21T09:47:43.728+0000"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sortAndSetFunctionList(tt.args.regionList, tt.args.runtimeList, tt.args.functionMap); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sortAndSetFunctionList() = %v, want %v", got, tt.want)
			}
		})
	}
}
