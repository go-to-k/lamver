package app

import (
	"context"
	"os"

	"github.com/go-to-k/lamver/internal/action"
	"github.com/go-to-k/lamver/internal/io"
	"github.com/go-to-k/lamver/internal/types"
	"github.com/go-to-k/lamver/pkg/client"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/urfave/cli/v2"
)

const SDKRetryMaxAttempts = 3

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
		Usage: "CLI tool to search Lambda runtime and versions.",
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

func (a *App) Run(ctx context.Context) error {
	return a.Cli.RunContext(ctx, os.Args)
}

func (a *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg, err := client.LoadAWSConfig(c.Context, a.DefaultRegion, a.Profile)
		if err != nil {
			return err
		}

		lambdaClient := client.NewLambda(
			lambda.NewFromConfig(cfg, func(o *lambda.Options) {
				o.RetryMaxAttempts = SDKRetryMaxAttempts
				o.RetryMode = aws.RetryModeStandard
			}),
		)

		ec2Client := client.NewEC2(
			ec2.NewFromConfig(cfg, func(o *ec2.Options) {
				o.RetryMaxAttempts = SDKRetryMaxAttempts
				o.RetryMode = aws.RetryModeStandard
			}),
		)

		getAllRegionsAndRuntimeInput := &action.GetAllRegionsAndRuntimeInput{
			Ctx:           c.Context,
			EC2:           ec2Client,
			Lambda:        lambdaClient,
			DefaultRegion: a.DefaultRegion,
		}
		allRegions, allRuntime, err := action.GetAllRegionsAndRuntime(getAllRegionsAndRuntimeInput)
		if err != nil {
			return err
		}

		regionsLabel := "Select regions you want to search.\n"
		targetRegions, continuation := io.GetCheckboxes(regionsLabel, allRegions)
		if !continuation {
			return nil
		}

		runtimeLabel := "Select runtime values you want to search.\n"
		targetRuntime, continuation := io.GetCheckboxes(runtimeLabel, allRuntime)
		if !continuation {
			return nil
		}

		keywordLabel := "Filter a keyword of function names(case-insensitive): "
		keyword := io.InputKeywordForFilter(keywordLabel)

		createFunctionListInput := &action.CreateFunctionListInput{
			Ctx:           c.Context,
			TargetRegions: targetRegions,
			TargetRuntime: targetRuntime,
			Keyword:       keyword,
			Lambda:        lambdaClient,
		}
		functionList, err := action.CreateFunctionList(createFunctionListInput)
		if err != nil {
			return err
		}

		functionHeader := types.GetLambdaFunctionDataKeys()
		if err := io.OutputResult(functionHeader, functionList, a.CSVOutputFilePath); err != nil {
			return err
		}

		return nil
	}
}
