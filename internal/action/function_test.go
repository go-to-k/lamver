package action

import (
	"context"
	"lamver/internal/types"
	"reflect"
	"testing"
)

func TestGetAllRegionsAndRuntime(t *testing.T) {
	type args struct {
		input *GetAllRegionsAndRuntimeInput
	}
	tests := []struct {
		name            string
		args            args
		wantRegionList  []string
		wantRuntimeList []string
		wantErr         bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRegionList, gotRuntimeList, err := GetAllRegionsAndRuntime(tt.args.input)
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
		ctx           context.Context
		region        string
		profile       string
		targetRuntime []string
		keyword       string
		functionCh    chan *types.LambdaFunctionData
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
			if err := putToFunctionChannelByRegion(tt.args.ctx, tt.args.region, tt.args.profile, tt.args.targetRuntime, tt.args.keyword, tt.args.functionCh); (err != nil) != tt.wantErr {
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
