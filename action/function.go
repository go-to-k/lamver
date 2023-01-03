package app

import (
	"context"
	"lamver/internal/types"
	"lamver/pkg/client"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

func GetAllRegionsAndRuntime(context context.Context, defaultRegion string, profile string) (regionList []string, runtimeList []string, err error) {
	cfg, err := loadAwsConfig(context, defaultRegion, profile)
	if err != nil {
		return regionList, runtimeList, err
	}

	eg, _ := errgroup.WithContext(context)

	eg.Go(func() error {
		ec2 := client.NewEC2(cfg)
		regionList, err = ec2.DescribeRegions(context)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		lambda := client.NewLambda(cfg)
		runtimeList = lambda.ListRuntimeValues()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return regionList, runtimeList, err
	}

	return regionList, runtimeList, nil
}

func CreateFunctionMap(ctx context.Context, profile string, targetRegions []string, targetRuntime []string, keyword string) (map[string]map[string][][]string, error) {
	functionMap := make(map[string]map[string][][]string, len(targetRuntime))

	eg, _ := errgroup.WithContext(ctx)
	wg := sync.WaitGroup{}
	functionCh := make(chan *types.LambdaFunctionData)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for f := range functionCh {
			if _, exist := functionMap[f.Runtime]; !exist {
				functionMap[f.Runtime] = make(map[string][][]string, len(targetRegions))
			}
			functionMap[f.Runtime][f.Region] = append(functionMap[f.Runtime][f.Region], []string{f.FunctionName, f.LastModified})
		}
	}()

	for _, region := range targetRegions {
		region := region
		eg.Go(func() error {
			return putToFunctionChannelByRegion(ctx, region, profile, targetRuntime, keyword, functionCh)
		})
	}

	go func() {
		eg.Wait()
		close(functionCh)
	}()

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	wg.Wait() // for functionMap race

	return functionMap, nil
}

func putToFunctionChannelByRegion(
	ctx context.Context,
	region string,
	profile string,
	targetRuntime []string,
	keyword string,
	functionCh chan *types.LambdaFunctionData,
) error {
	cfg, err := loadAwsConfig(ctx, region, profile)
	if err != nil {
		return err
	}

	lambda := client.NewLambda(cfg)
	functions, err := lambda.ListFunctions(ctx)
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

func SortAndSetFunctionList(regionList []string, runtimeList []string, functionMap map[string]map[string][][]string) [][]string {
	var functionData [][]string

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
				functionData = append(functionData, data)
			}
		}
	}

	return functionData
}
