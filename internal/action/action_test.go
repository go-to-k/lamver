package action

import (
	"context"
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
				m.EXPECT().DescribeRegions(ctx).Return(
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

func TestCreateFunctionMap(t *testing.T) {
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
		want                      map[string]map[string][][]string
		wantErr                   bool
	}{
		{
			name: "CreateFunctionMap success",
			args: args{
				ctx:           ctx,
				targetRegions: []string{"ap-northeast-1", "us-east-1", "us-east-2"},
				targetRuntime: []string{"nodejs18.x", "nodejs"},
				keyword:       "",
			},
			prepareMockLambdaClientFn: func(m *client.MockLambdaClient) {
				m.EXPECT().ListFunctionsWithRegion(ctx, "ap-northeast-1").Return(
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
				m.EXPECT().ListFunctionsWithRegion(ctx, "us-east-1").Return(
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
				m.EXPECT().ListFunctionsWithRegion(ctx, "us-east-2").Return(
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
			want: map[string]map[string][][]string{
				"nodejs18.x": {
					"us-east-2": {
						[]string{"function6", "2022-12-22T09:47:43.728+0000"},
					},
				},
				"nodejs": {
					"ap-northeast-1": {
						[]string{"function1", "2022-12-21T09:47:43.728+0000"},
					},
					"us-east-1": {
						[]string{"function3", "2022-12-21T09:47:43.728+0000"},
					},
					"us-east-2": {
						[]string{"function5", "2022-12-21T09:47:43.728+0000"},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			lambdaClientMock := client.NewMockLambdaClient(ctrl)

			tt.prepareMockLambdaClientFn(lambdaClientMock)

			input := &CreateFunctionMapInput{
				Ctx:           tt.args.ctx,
				TargetRegions: tt.args.targetRegions,
				TargetRuntime: tt.args.targetRuntime,
				Keyword:       tt.args.keyword,
				Lambda:        lambdaClientMock,
			}

			got, err := CreateFunctionMap(input)
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
				m.EXPECT().ListFunctionsWithRegion(ctx, "us-east-1").Return(
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

			go func() {
				for {
					select {
					case <-ctx.Done():
						return
					case <-tt.args.functionCh:
					default:
					}
				}
			}()

			if err := putToFunctionChannelByRegion(tt.args.ctx, tt.args.region, tt.args.targetRuntime, tt.args.keyword, tt.args.functionCh, lambdaClientMock); (err != nil) != tt.wantErr {
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
		{
			name: "SortAndSetFunctionList success",
			args: args{
				input: &SortAndSetFunctionListInput{
					RegionList:  []string{"ap-northeast-1", "us-east-1", "us-east-2"},
					RuntimeList: []string{"nodejs18.x", "nodejs"},
					FunctionMap: map[string]map[string][][]string{
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
			if got := SortAndSetFunctionList(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SortAndSetFunctionList() = %v, want %v", got, tt.want)
			}
		})
	}
}
