package main

import (
	"context"
	"lamver/internal/app"
	"lamver/internal/logger"
	"lamver/internal/version"
	"os"
)

func main() {
	logger.NewLogger(version.IsDebug())
	ctx := context.Background()
	app := app.NewApp(version.GetVersion())

	if err := app.Run(ctx); err != nil {
		logger.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}
