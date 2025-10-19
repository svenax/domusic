package cmd

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/urfave/cli/v3"
)

var cfgFile string

// Execute creates and runs the CLI application
func Execute() {
	app := &cli.Command{
		Name:  "domusic",
		Usage: "Handles the music library at svenax.net",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Usage:       getConfigUsage(),
				Destination: &cfgFile,
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
			versionCmd,
			viewCmd,
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func getConfigUsage() string {
	filename := ".domusic.yaml"
	if runtime.GOOS == "windows" {
		filename = "domusic.ini"
	}
	return fmt.Sprintf("config file (default is $HOME/%s)", filename)
}
