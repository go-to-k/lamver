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

// TODO: CSV files and JSON output option
// TODO: refactor for nested loops
// TODO: write app tests for regions
// TODO: write sdk tests not using interface, otherwise use interface, go mock and auto creating test modules

func (app *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		functionHeader := types.GetLambdaFunctionDataKeys()
		functionData := [][]string{}

		cfg, err := app.loadAwsConfig(c.Context, app.DefaultRegion)
		if err != nil {
			return err
		}

		eg, _ := errgroup.WithContext(c.Context)
		regionsLabel := "Select regions you want to display.\n"
		runtimeLabel := "Select runtime values you want to display.\n"

		var (
			regionList  []string
			runtimeList []string
		)

		eg.Go(func() error {
			ec2 := client.NewEC2Client(cfg)
			regionList, err = ec2.DescribeRegions(c.Context)
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
			return err
		}

		targetRegions, continuation := io.GetCheckboxes(regionsLabel, regionList)
		if !continuation {
			return nil
		}

		targetRuntime, continuation := io.GetCheckboxes(runtimeLabel, runtimeList)
		if !continuation {
			return nil
		}

		keyword := io.InputKeywordForFilter()

		eg, _ = errgroup.WithContext(c.Context)
		wg := sync.WaitGroup{}
		functionMap := make(map[string]map[string][][]string, len(targetRuntime))
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
				cfg, err := app.loadAwsConfig(c.Context, region)
				if err != nil {
					return err
				}

				lambda := client.NewLambdaClient(cfg)
				functions, err := lambda.ListFunctions(c.Context)
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
			})
		}

		go func() {
			eg.Wait()
			close(functionCh)
		}()

		if err := eg.Wait(); err != nil {
			return err
		}

		wg.Wait() // for functionMap race

		// sort and set to functionData
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

		if err := io.OutputResult(functionHeader, functionData, app.CSVOutputFilePath); err != nil {
			return err
		}

		return nil
	}
}
