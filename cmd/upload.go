package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/grokify/html-strip-tags-go"
	"github.com/spf13/cobra"
)

var (
	verbose bool
)

// Upload a pdf file to Evernote and create/update a note with the pdf
// attached.
//
// NOTE: This code depends on geeknote for Evernote integration.
//       It should move to using the Evernote API directly at some point.
var uploadCmd = &cobra.Command{
	Use:   "upload file",
	Short: "Upload and/or update file to Evernote",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			files := []string{arg}
			if strings.Contains(arg, "*") {
				files, _ = filepath.Glob(pathFromRoot(arg))
			}
			for _, f := range files {
				f = getOutputPath(f, false)
				if _, err := os.Stat(f); err != nil {
					errExit(err)
				}
				run(f)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)

	uploadCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "show all command output")
}

func run(file string) {
	title := humanize(strings.TrimSuffix(path.Base(file), path.Ext(file)))
	tag := tagify(path.Base(path.Dir(file)))
	notebook, err := getNotebook()
	errExit(err)

	out, err := upload("find", "--exact-entry", "--guid", "--search", title,
		"--notebook", notebook)
	outlines := strings.Split(strings.Trim(out, "\n"), "\n")
	outlast := outlines[len(outlines)-1]

	if err != nil {
		if strings.Contains(outlast, "Rate Limit") {
			msgExit(outlast)
		}

		fmt.Println("Creating:", title, "in", tag)
		out, err = upload("create", "--notebook", notebook, "--title", title,
			"--content", "", "--resource", file, "--tag", tag)
		if err != nil {
			msgExit(out)
		}
	} else {
		fmt.Println("Updating:", title, "in", tag)
		guid := strings.Split(outlast, " ")[0]
		out, err := upload("show", "--note", guid)
		if err != nil {
			fmt.Println(err)
		}

		rx := regexp.MustCompile("(?s)CONTENT -+\nTags: [^\n]*\n(.*)")
		content := rx.FindStringSubmatch(out)[1]
		if content != "" {
			content = strings.Trim(strip.StripTags(content), " \n")
		}

		_, err = upload("edit", "--note", guid,
			"--content", content, "--resource", file)
		errExit(err)
	}
}

func upload(args ...string) (string, error) {
	if verbose {
		fmt.Println("CMD: geeknote", strings.Join(args, " "))
	}

	c := exec.Command("geeknote", args...)
	out, err := c.CombinedOutput()
	outstr := string(out)

	if verbose {
		fmt.Println("RES:", outstr)
	}

	return outstr, err
}

func humanize(text string) string {
	text = strings.TrimPrefix(text, "!")
	text = strings.Replace(text, "_", " ", -1)
	text = strings.Title(text)

	return text
}

func tagify(tag string) string {
	tag = strings.TrimPrefix(tag, "!")
	tag = strings.Replace(tag, "-", "/", -1)
	tag = strings.Replace(tag, "_", " ", -1)
	if strings.HasSuffix(tag, "pes") || strings.HasSuffix(tag, "tes") {
		tag = strings.TrimSuffix(tag, "s")
	}
	tag = strings.TrimSuffix(tag, "es")
	tag = strings.TrimSuffix(tag, "s")
	tag = strings.Title(tag)

	return tag
}
