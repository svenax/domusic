package cmd

import (
	"os/exec"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit file",
	Short: "Create or edit a Lilypond music file",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			errExit("edit needs a file name")
		}

		e, ea, err := getEditor()
		if err != nil {
			errExit(err)
		}

		c := exec.Command(e, append(ea, getSourcePath(args[0]))...)
		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
