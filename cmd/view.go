package cmd

import (
	"fmt"
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
			if os.IsNotExist(err) {
				msgExit(fmt.Sprintf("output file does not exist: %s", file))
			}
			errExit(err)
		}

		v, _, err := getViewer()
		errExit(err)
		c := exec.Command("open", "-b", v, file)
		if err := c.Run(); err != nil {
			errExit(fmt.Errorf("failed to open file with viewer '%s' for file '%s': %w", v, file, err))
		}
	},
}

func init() {
	rootCmd.AddCommand(viewCmd)

	viewCmd.Flags().BoolVarP(&preview, "preview", "p", false, "view the preview image")
}
