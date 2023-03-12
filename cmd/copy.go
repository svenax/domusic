package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var lyCopier *copier
var copyCmd = &cobra.Command{
	Use:   "copy files...",
	Short: "Copy output files, at the same time renaming them for Forscore",
	Long: `Copy output files, at the same time renaming them for Forscore.

The --output path must be set and the folder exist, otherwise the command will fail.
All files will be copied to a flat structure inside output path. This is how they are
stored in Forscore.`,
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			files := []string{arg}
			if strings.Contains(arg, "*") {
				files, _ = filepath.Glob(pathFromRoot(arg))
			}
			for _, f := range files {
				lyCopier.copy(getSourcePath(f))
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)

	copyCmd.Flags().BoolP("dry-run", "d", false, "only show copy commands")
	copyCmd.Flags().StringP("output", "o", "", "output directory (required)")
	copyCmd.MarkFlagRequired("output")

	lyCopier = &copier{copyCmd}
}

type copier struct {
	cmd *cobra.Command
}

func (c *copier) copy(src string) {
	source, err := os.ReadFile(src)
	errExit(err)

	matches := titleRx.FindSubmatch(source)
	if matches == nil {
		log.Printf("Skipping '%s'. No title field found", src)
	} else {
		target := filepath.Join(c.flagString("output"), string(matches[1])+".pdf")
		if c.flagBool("dry-run") {
			fmt.Printf("cp %s '%s'\n", getPdfPath(src), target)
		} else {
			_, err = copyFile(getPdfPath(src), target)
			errExit(err)
		}
	}
}

func (c *copier) flagBool(name string) bool {
	val, _ := c.cmd.Flags().GetBool(name)
	return val
}

func (c *copier) flagString(name string) string {
	val, _ := c.cmd.Flags().GetString(name)
	return val
}
