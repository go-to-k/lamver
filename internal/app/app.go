package app

import (
	"context"
	"lamver/internal/action"
	"lamver/internal/io"
	"lamver/internal/types"
	"lamver/pkg/client"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/urfave/cli/v2"
)

const awsSDKRetryMaxAttempts = 3

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
				o.RetryMaxAttempts = awsSDKRetryMaxAttempts
				o.RetryMode = aws.RetryModeStandard
			}),
		)

		ec2Client := client.NewEC2(
			ec2.NewFromConfig(cfg, func(o *ec2.Options) {
				o.RetryMaxAttempts = awsSDKRetryMaxAttempts
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

		regionsLabel := "Select regions you want to display.\n"
		targetRegions, continuation := io.GetCheckboxes(regionsLabel, allRegions)
		if !continuation {
			return nil
		}

		runtimeLabel := "Select runtime values you want to display.\n"
		targetRuntime, continuation := io.GetCheckboxes(runtimeLabel, allRuntime)
		if !continuation {
			return nil
		}

		keyword := io.InputKeywordForFilter()

		createFunctionMapInput := &action.CreateFunctionListInput{
			Ctx:           c.Context,
			TargetRegions: targetRegions,
			TargetRuntime: targetRuntime,
			Keyword:       keyword,
			Lambda:        lambdaClient,
		}
		functionData, err := action.CreateFunctionList(createFunctionMapInput)
		if err != nil {
			return err
		}

		functionHeader := types.GetLambdaFunctionDataKeys()
		if err := io.OutputResult(functionHeader, functionData, a.CSVOutputFilePath); err != nil {
			return err
		}

		return nil
	}
}
