package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit file",
	Short: "Create or edit a Lilypond music file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			msgExit("edit needs a file name")
		}

		e, ea, err := getEditor()
		errExit(err)

		c := exec.Command(e, append(ea, getSourcePath(args[0]))...)
		if err := c.Run(); err != nil {
			errExit(fmt.Errorf("failed to run editor '%s %v': %w", e, append(ea, getSourcePath(args[0])), err))
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
