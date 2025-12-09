package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/urfave/cli/v3"
)

var viewCmd = &cli.Command{
	Name:  "view",
	Usage: "View PDF or preview image <file>",
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:    "preview",
			Aliases: []string{"p"},
			Usage:   "view the preview image",
		},
	},
	Arguments: []cli.Argument{
		&cli.StringArg{
			Name: "file",
		},
	},
	Action: func(ctx context.Context, cmd *cli.Command) error {
		preview := cmd.Bool("preview")
		file := getOutputPath(cmd.StringArg("file"), preview)
		if _, err := os.Stat(file); err != nil {
			if os.IsNotExist(err) {
				return printAndReturnError("output file does not exist: %s", file)
			}
			return printAndReturnError("failed to stat file: %w", err)
		}

		v, _, err := getViewer()
		if err != nil {
			return printAndReturnError("failed to get viewer: %w", err)
		}

		c := exec.Command("open", "-b", v, file)
		if err := c.Run(); err != nil {
			return printAndReturnError("failed to open file with viewer '%s' for file '%s': %w", v, file, err)
		}
		return nil
	},
}
