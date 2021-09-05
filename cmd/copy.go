package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy to_dir",
	Short: "Copy output files, at the same time renaming them for Forscore",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			files := []string{arg}
			if strings.Contains(arg, "*") {
				files, _ = filepath.Glob(pathFromRoot(arg))
			}
			target := "/akk/folder"
			for _, f := range files {
				source, err := ioutil.ReadFile(getSourcePath(f))
				errExit(err)

				re := regexp.MustCompile(`title\s*=\s*"(.*?)"`)
				matches := re.FindSubmatch(source)
				if matches == nil {
					log.Printf("Skipping '%s'. No title field found", getSourcePath(f))
				} else {
					fmt.Printf("cp %s %s/%s.pdf\n", getPdfPath(f), target, string(matches[1]))
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}
