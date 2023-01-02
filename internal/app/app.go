package app

import (
	"context"
	"lamver/internal/io"
	"lamver/internal/types"
	"lamver/pkg/client"
	"os"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Cli               *cli.App
	Profile           string
	DefaultRegion     string
	CSVOutputFilePath string
}

func NewApp(version string) *App {
	app := App{}

	app.Cli = &cli.App{
		Name:  "lamver",
		Usage: "CLI tool to display Lambda runtime and versions.",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "profile",
				Aliases:     []string{"p"},
				Usage:       "AWS profile name",
				Destination: &app.Profile,
			},
			&cli.StringFlag{
				Name:        "region",
				Aliases:     []string{"r"},
				Usage:       "AWS default region",
				Destination: &app.DefaultRegion,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "Output file path for CSV format",
				Destination: &app.CSVOutputFilePath,
			},
		},
	}

	app.Cli.Version = version
	app.Cli.Action = app.getAction()
	app.Cli.HideHelpCommand = true

	return &app
}

func (app *App) Run(ctx context.Context) error {
	return app.Cli.RunContext(ctx, os.Args)
}

func (app *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		regionList, runtimeList, err := app.getRuntimeAndRegions(c.Context)
		if err != nil {
			return err
		}

		regionsLabel := "Select regions you want to display.\n"
		targetRegions, continuation := io.GetCheckboxes(regionsLabel, regionList)
		if !continuation {
			return nil
		}

		runtimeLabel := "Select runtime values you want to display.\n"
		targetRuntime, continuation := io.GetCheckboxes(runtimeLabel, runtimeList)
		if !continuation {
			return nil
		}

		keyword := io.InputKeywordForFilter()

		functionMap, err := app.createFunctionMap(c.Context, targetRegions, targetRuntime, keyword)
		if err != nil {
			return err
		}

		functionHeader := types.GetLambdaFunctionDataKeys()
		functionData := app.sortFunctionMap(regionList, runtimeList, functionMap)
		if err := io.OutputResult(functionHeader, functionData, app.CSVOutputFilePath); err != nil {
			return err
		}

		return nil
	}
}

func (app *App) getRuntimeAndRegions(context context.Context) (regionList []string, runtimeList []string, err error) {
	cfg, err := app.loadAwsConfig(context, app.DefaultRegion)
	if err != nil {
		return regionList, runtimeList, err
	}

	eg, _ := errgroup.WithContext(context)

	eg.Go(func() error {
		ec2 := client.NewEC2Client(cfg)
		regionList, err = ec2.DescribeRegions(context)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		lambda := client.NewLambdaClient(cfg)
		runtimeList = lambda.ListRuntimeValues()
		return nil
	})

	if err := eg.Wait(); err != nil {
		return regionList, runtimeList, err
	}

	return regionList, runtimeList, nil
}

func (app *App) createFunctionMap(ctx context.Context, targetRegions []string, targetRuntime []string, keyword string) (map[string]map[string][][]string, error) {
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
			return app.putToLambdaFunctionChannelByRegion(ctx, region, targetRuntime, keyword, functionCh)
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

func (app *App) putToLambdaFunctionChannelByRegion(
	ctx context.Context,
	region string,
	targetRuntime []string,
	keyword string,
	functionCh chan *types.LambdaFunctionData,
) error {
	cfg, err := app.loadAwsConfig(ctx, region)
	if err != nil {
		return err
	}

	lambda := client.NewLambdaClient(cfg)
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

func (app *App) sortFunctionMap(regionList []string, runtimeList []string, functionMap map[string]map[string][][]string) [][]string {
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
