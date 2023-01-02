package app

import (
	"context"
	"lamver/internal/action"
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

func (a *App) Run(ctx context.Context) error {
	return a.Cli.RunContext(ctx, os.Args)
}

func (a *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		getAllRegionsAndRuntimeInput := &action.GetAllRegionsAndRuntimeInput{
			Ctx:           c.Context,
			DefaultRegion: a.DefaultRegion,
			Profile:       a.Profile,
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

		createFunctionMapInput := &action.CreateFunctionMapInput{
			Ctx:           c.Context,
			Profile:       a.Profile,
			TargetRegions: targetRegions,
			TargetRuntime: targetRuntime,
			Keyword:       keyword,
		}
		functionMap, err := action.CreateFunctionMap(createFunctionMapInput)
		if err != nil {
			return err
		}

		sortAndSetFunctionListInput := &action.SortAndSetFunctionListInput{
			RegionList:  targetRegions,
			RuntimeList: targetRuntime,
			FunctionMap: functionMap,
		}
		functionData := action.SortAndSetFunctionList(sortAndSetFunctionListInput)

		functionHeader := types.GetLambdaFunctionDataKeys()
		if err := io.OutputResult(functionHeader, functionData, a.CSVOutputFilePath); err != nil {
			return err
		}

		return nil
	}
}
