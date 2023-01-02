package app

import (
	"bufio"
	"context"
	"fmt"
	"lamver/internal/logger"
	"lamver/internal/types"
	"lamver/pkg/client"
	"os"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

const DEFAULT_AWS_REGION = "us-east-1"

type App struct {
	Cli           *cli.App
	Profile       string
	DefaultRegion string
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

// TODO: selection(filter) of keyword(prefix, suffix, etc...) of function names
// TODO: max items(input to params for lambda methods)
// TODO: sort (order by count desc)?
// TODO: write app tests for regions
// TODO: write sdk tests not using interface, otherwise use interface, go mock and auto creating test modules
// TODO: aggregate output option
// TODO: CSV files and JSON output option

func (app *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {

		functionHeader := []string{"Runtime", "Region", "FunctionName", "LastModified"}
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

		targetRegions, continuation := getCheckboxes(regionsLabel, regionList)
		if !continuation {
			return nil
		}

		targetRuntime, continuation := getCheckboxes(runtimeLabel, runtimeList)
		if !continuation {
			return nil
		}

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

		// TODO: refactor for nested loops
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
						if string(function.Runtime) == runtime {
							functionCh <- &types.LambdaFunctionData{
								Runtime:      runtime,
								Region:       region,
								FunctionName: *function.FunctionName,
								LastModified: *function.LastModified,
							}
							break
						}
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

		fmt.Println(*logger.ToStringAsTableFormat(functionHeader, functionData))
		return nil
	}
}

func (app *App) loadAwsConfig(ctx context.Context, region string) (aws.Config, error) {
	var (
		cfg aws.Config
		err error
	)

	if app.Profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(app.Profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
	}

	if region != "" {
		cfg.Region = region
	}
	if cfg.Region == "" {
		cfg.Region = DEFAULT_AWS_REGION
	}

	return cfg, err
}

func getCheckboxes(label string, opts []string) ([]string, bool) {
	var checkboxes []string

	for {
		prompt := &survey.MultiSelect{
			Message: label,
			Options: opts,
		}
		survey.AskOne(prompt, &checkboxes)

		if len(checkboxes) == 0 {
			logger.Logger.Warn().Msg("Select values!")
			ok := getYesNo("Do you want to finish?")
			if ok {
				logger.Logger.Info().Msg("Finished...")
				return checkboxes, false
			}
			continue
		}

		ok := getYesNo("OK?")
		if ok {
			return checkboxes, true
		}
	}
}

func getYesNo(label string) bool {
	choices := "Y/n"
	r := bufio.NewReader(os.Stdin)
	var s string

	for {
		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
		s, _ = r.ReadString('\n')
		fmt.Fprintln(os.Stderr)

		s = strings.TrimSpace(s)
		if s == "" {
			return true
		}
		s = strings.ToLower(s)
		if s == "y" || s == "yes" {
			return true
		}
		if s == "n" || s == "no" {
			return false
		}
	}
}
