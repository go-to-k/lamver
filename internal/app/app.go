package app

import (
	"context"
	"lamver/internal/io"
	"lamver/internal/types"
	"os"

	"github.com/urfave/cli/v2"
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
