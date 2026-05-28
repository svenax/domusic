package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v3"
)

var editCmd = &cli.Command{
	Name:  "edit",
	Usage: "Create or edit a Lilypond music file <file>",
	Action: func(ctx context.Context, cmd *cli.Command) error {
		args := cmd.Args().Slice()
		if len(args) != 1 {
			return fmt.Errorf("edit needs a file name")
		}

		e, ea, err := getEditor()
		if err != nil {
			return fmt.Errorf("failed to get editor: %w", err)
		}

		c := exec.Command(e, append(ea, getSourcePath(args[0]))...)
		if err := c.Run(); err != nil {
			return fmt.Errorf("failed to run editor '%s %v': %w", e, append(ea, getSourcePath(args[0])), err)
		}
		return nil
	},
}
