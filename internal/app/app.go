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

type App struct {
	Cli             *cli.App
	Profile         string
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

func (app *App) getAction() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		config, err := app.loadAwsConfig(c.Context)
		if err != nil {
			return err
		}

		ec2 := client.NewEC2Client(config)
		regions, err := ec2.DescribeRegions(c.Context)
		if err != nil {
			return err
		}

		for _, region := range regions {
			fmt.Println(region)
		}
		return nil
	}
}

func (app *App) loadAwsConfig(ctx context.Context) (aws.Config, error) {
	var (
		cfg aws.Config
		err error
	)

	if app.Profile != "" {
		// cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(app.Region), config.WithSharedConfigProfile(app.Profile))
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(app.Profile))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx)
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
