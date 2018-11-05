package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var preview bool

var viewCmd = &cobra.Command{
	Use:   "view file",
	Short: "View PDF or preview image",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			msgExit("view needs a file name")
		}

		file := getOutputPath(args[0], preview)
		if _, err := os.Stat(file); err != nil {
			errExit(err)
		}

		v, _, err := getViewer()
		errExit(err)
		c := exec.Command("open", "-b", v, file)
		c.Run()
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVarP(&preview, "preview", "p", false, "view the preview image")
}
