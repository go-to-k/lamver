package app

import (
	"context"
	"fmt"
	"lamver/pkg/client"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/urfave/cli/v2"
)

const DEFAULT_AWS_REGION = "us-east-1"

type App struct {
	Cli             *cli.App
	Profile         string
	DefaultRegion   string
	InteractiveMode bool
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
			&cli.BoolFlag{
				Name:        "interactive",
				Aliases:     []string{"i"},
				Value:       false,
				Usage:       "Interactive Mode",
				Destination: &app.InteractiveMode,
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

// TODO: display for table style
// TODO: interactive modes
// TODO: selection(filter) of kind of Lambda runtime
// TODO: selection(filter) of regions
// TODO: selection(filter) of keyword(prefix, suffix, etc...) of function names
// TODO: max items(input to params for lambda methods)
// TODO: sort (order by count desc)?
// TODO: write app tests for regions
// TODO: write sdk tests not using interface, otherwise use interface, go mock and auto creating test modules
func (app *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		cfg, err := app.loadAwsConfig(c.Context, app.DefaultRegion)
		if err != nil {
			return err
		}

		ec2 := client.NewEC2Client(cfg)
		regions, err := ec2.DescribeRegions(c.Context)
		if err != nil {
			return err
		}
		for _, region := range regions {
			fmt.Println(region)
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
				fmt.Println(*function.FunctionName, function.Runtime)
			}
			fmt.Println(len(functions))
		}

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

// func doInteractiveMode() ([]string, bool) {
// 	var checkboxes []string

// 	for {
// 		checkboxes = getCheckboxes()

// 		if len(checkboxes) == 0 {
// 			logger.Logger.Warn().Msg("Select ResourceTypes!")
// 			ok := getYesNo("Do you want to finish?")
// 			if ok {
// 				logger.Logger.Info().Msg("Finished...")
// 				return checkboxes, false
// 			}
// 			continue
// 		}

// 		ok := getYesNo("OK?")
// 		if ok {
// 			return checkboxes, true
// 		}
// 	}
// }

// func getCheckboxes() []string {
// 	label := "Select ResourceTypes you wish to delete even if DELETE_FAILED." +
// 		"\n" +
// 		"However, if resources of the selected ResourceTypes will not be DELETE_FAILED when the stack is deleted, the resources will be deleted even if you selected. " +
// 		"\n"
// 	opts := resourcetype.GetResourceTypes()
// 	res := []string{}

// 	prompt := &survey.MultiSelect{
// 		Message: label,
// 		Options: opts,
// 	}
// 	survey.AskOne(prompt, &res)

// 	return res
// }

// func getYesNo(label string) bool {
// 	choices := "Y/n"
// 	r := bufio.NewReader(os.Stdin)
// 	var s string

// 	for {
// 		fmt.Fprintf(os.Stderr, "%s (%s) ", label, choices)
// 		s, _ = r.ReadString('\n')
// 		fmt.Fprintln(os.Stderr)

// 		s = strings.TrimSpace(s)
// 		if s == "" {
// 			return true
// 		}
// 		s = strings.ToLower(s)
// 		if s == "y" || s == "yes" {
// 			return true
// 		}
// 		if s == "n" || s == "no" {
// 			return false
// 		}
// 	}
// }
