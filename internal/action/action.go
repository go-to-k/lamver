package action

import (
	"context"
	"lamver/internal/types"
	"lamver/pkg/client"
	"strings"

	"golang.org/x/sync/errgroup"
)

type GetAllRegionsAndRuntimeInput struct {
	Ctx           context.Context
	EC2           client.EC2Client
	Lambda        client.LambdaClient
	DefaultRegion string
}

func GetAllRegionsAndRuntime(input *GetAllRegionsAndRuntimeInput) (regionList []string, runtimeList []string, err error) {
	eg, ctx := errgroup.WithContext(input.Ctx)
	eg.Go(func() error {
		regionList, err = input.EC2.DescribeRegions(ctx)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		runtimeList = input.Lambda.ListRuntimeValues()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return regionList, runtimeList, err
	}

	return regionList, runtimeList, nil
}

type CreateFunctionListInput struct {
	Ctx           context.Context
	TargetRegions []string
	TargetRuntime []string
	Keyword       string
	Lambda        client.LambdaClient
}

func CreateFunctionList(input *CreateFunctionListInput) ([][]string, error) {
	functionMap := make(map[string]map[string][][]string, len(input.TargetRuntime))

	eg, ctx := errgroup.WithContext(input.Ctx)
	functionCh := make(chan *types.LambdaFunctionData)

	for _, region := range input.TargetRegions {
		region := region
		eg.Go(func() error {
			return putToFunctionChannelByRegion(
				ctx,
				region,
				input.TargetRuntime,
				input.Keyword,
				functionCh,
				input.Lambda,
			)
		})
	}

	go func() {
		eg.Wait()
		close(functionCh)
	}()

	for f := range functionCh {
		if _, exist := functionMap[f.Runtime]; !exist {
			functionMap[f.Runtime] = make(map[string][][]string, len(input.TargetRegions))
		}
		functionMap[f.Runtime][f.Region] = append(functionMap[f.Runtime][f.Region], []string{f.FunctionName, f.LastModified})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	sortedFunctionList := sortAndSetFunctionList(input.TargetRegions, input.TargetRuntime, functionMap)

	return sortedFunctionList, nil
}

func putToFunctionChannelByRegion(
	ctx context.Context,
	region string,
	targetRuntime []string,
	keyword string,
	functionCh chan *types.LambdaFunctionData,
	lambda client.LambdaClient,
) error {
	functions, err := lambda.ListFunctionsWithRegion(ctx, region)
	if err != nil {
		return err
	}

	for _, function := range functions {
		for _, runtime := range targetRuntime {
			if string(function.Runtime) != runtime {
				continue
			}
			if strings.Contains(*function.FunctionName, keyword) {
				functionCh <- &types.LambdaFunctionData{
					Runtime:      runtime,
					Region:       region,
					FunctionName: *function.FunctionName,
					LastModified: *function.LastModified,
				}
			}
			break
		}
	}

	return nil
}

func sortAndSetFunctionList(regionList []string, runtimeList []string, functionMap map[string]map[string][][]string) [][]string {
	var functionList [][]string

	for _, runtime := range runtimeList {
		if _, exist := functionMap[runtime]; !exist {
			continue
		}
		for _, region := range regionList {
			if _, exist := functionMap[runtime][region]; !exist {
				continue
			}
			for _, f := range functionMap[runtime][region] {
				var data []string
				data = append(data, runtime)
				data = append(data, region)
				data = append(data, f...)
				functionList = append(functionList, data)
			}
		}
	}

	return functionList
}
