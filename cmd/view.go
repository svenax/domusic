package cmd

import (
	"context"
	"fmt"
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
				cli.Exit(fmt.Sprintf("output file does not exist: %s", file), 1)
			}
			cli.Exit(err, 1)
		}

		if v, _, err := getViewer(); err != nil {
			cli.Exit(err, 1)
		} else {
			c := exec.Command("open", "-b", v, file)
			if err := c.Run(); err != nil {
				cli.Exit(fmt.Errorf("failed to open file with viewer '%s' for file '%s': %w", v, file, err), 1)
			}
		}
		return nil
	},
}
