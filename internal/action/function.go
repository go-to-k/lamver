package action

import (
	"context"
	"lamver/internal/types"
	"lamver/pkg/client"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"
)

type GetAllRegionsAndRuntimeInput struct {
	Ctx              context.Context
	AWSConfigFactory *AWSConfigFactory
	EC2Factory       *client.EC2Factory
	LambdaFactory    *client.LambdaFactory
	DefaultRegion    string
}

func GetAllRegionsAndRuntime(input *GetAllRegionsAndRuntimeInput) (regionList []string, runtimeList []string, err error) {
	cfg, err := input.AWSConfigFactory.Create(input.Ctx, input.DefaultRegion)
	if err != nil {
		return regionList, runtimeList, err
	}

	eg, _ := errgroup.WithContext(input.Ctx)
	eg.Go(func() error {
		ec2 := input.EC2Factory.Create(cfg.Config)
		regionList, err = ec2.DescribeRegions(input.Ctx)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		lambda := input.LambdaFactory.Create(cfg.Config)
		runtimeList = lambda.ListRuntimeValues()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return regionList, runtimeList, err
	}

	return regionList, runtimeList, nil
}

type CreateFunctionMapInput struct {
	Ctx              context.Context
	TargetRegions    []string
	TargetRuntime    []string
	Keyword          string
	AWSConfigFactory *AWSConfigFactory
	LambdaFactory    *client.LambdaFactory
}

func CreateFunctionMap(input *CreateFunctionMapInput) (map[string]map[string][][]string, error) {
	functionMap := make(map[string]map[string][][]string, len(input.TargetRuntime))

	eg, _ := errgroup.WithContext(input.Ctx)
	wg := sync.WaitGroup{}
	functionCh := make(chan *types.LambdaFunctionData)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for f := range functionCh {
			if _, exist := functionMap[f.Runtime]; !exist {
				functionMap[f.Runtime] = make(map[string][][]string, len(input.TargetRegions))
			}
			functionMap[f.Runtime][f.Region] = append(functionMap[f.Runtime][f.Region], []string{f.FunctionName, f.LastModified})
		}
	}()

	for _, region := range input.TargetRegions {
		region := region
		eg.Go(func() error {
			return putToFunctionChannelByRegion(
				input.Ctx,
				region,
				input.TargetRuntime,
				input.Keyword,
				functionCh,
				input.AWSConfigFactory,
				input.LambdaFactory,
			)
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
	targetRuntime []string,
	keyword string,
	functionCh chan *types.LambdaFunctionData,
	awsConfigFactory *AWSConfigFactory,
	lambdaFactory *client.LambdaFactory,
) error {
	cfg, err := awsConfigFactory.Create(ctx, region)
	if err != nil {
		return err
	}

	lambda := lambdaFactory.Create(cfg.Config)
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

type SortAndSetFunctionListInput struct {
	RegionList  []string
	RuntimeList []string
	FunctionMap map[string]map[string][][]string
}

func SortAndSetFunctionList(input *SortAndSetFunctionListInput) [][]string {
	var functionData [][]string

	for _, runtime := range input.RuntimeList {
		if _, exist := input.FunctionMap[runtime]; !exist {
			continue
		}
		for _, region := range input.RegionList {
			if _, exist := input.FunctionMap[runtime][region]; !exist {
				continue
			}
			for _, f := range input.FunctionMap[runtime][region] {
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
