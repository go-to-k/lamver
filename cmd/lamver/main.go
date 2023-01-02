package main

import (
	"context"
	"lamver/internal/app"
	"lamver/internal/io"
	"lamver/internal/version"
	"os"
)

func main() {
	io.NewLogger(version.IsDebug())
	ctx := context.Background()
	app := app.NewApp(version.GetVersion())

	if err := app.Run(ctx); err != nil {
		io.Logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}
