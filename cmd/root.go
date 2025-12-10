package cmd

import (
	"context"
	"os"
	"runtime"

	"github.com/urfave/cli/v3"
)

var configPath string

// Execute creates and runs the CLI application
func Execute() {
	app := &cli.Command{
		Name:  "domusic",
		Usage: "Handles the music library at svenax.net",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       getConfigUsage(),
				Destination: &configPath,
			},
		},
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			initConfig()
			return ctx, nil
		},
		Commands: []*cli.Command{
			collectionCmd,
			editCmd,
			makeCmd,
			syncCmd,
			versionCmd,
			viewCmd,
		},
	}

	if app.Run(context.Background(), os.Args) != nil {
		// Error messages are already printed elsewhere
		os.Exit(1)
	}
}

func getConfigUsage() string {
	if runtime.GOOS == "windows" {
		return "config file (default searches: .\\.domusic.yaml, %APPDATA%\\domusic\\config.yaml, %HOME%\\.domusic.yaml)"
	}
	return "config file (default searches: ./.domusic.yaml (and parents), ~/.config/domusic/config.yaml, ~/.domusic.yaml)"
}
