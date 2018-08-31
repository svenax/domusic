package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var lyMaker *maker
var makeCmd = &cobra.Command{
	Use:   "make files...",
	Short: "Run Lilypond on music file(s)",
	Run: func(cmd *cobra.Command, args []string) {
		for _, arg := range args {
			files := []string{arg}
			if strings.Contains(arg, "*") {
				files, _ = filepath.Glob(pathFromRoot(arg))
			}
			for _, f := range files {
				lyMaker.run(getSourcePath(f))
			}
		}
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if p, _ := cmd.Flags().GetBool("post"); p {
			cmd.Flags().Set("type", "png")
			cmd.Flags().Set("root", "true")
			cmd.Flags().Set("crop", "true")
		}
	},
}

func init() {
	rootCmd.AddCommand(makeCmd)

	makeCmd.Flags().IntP("staff-size", "s", 15, "staff size")
	makeCmd.Flags().StringP("paper-size", "p", "a4", "paper size")
	makeCmd.Flags().StringP("format", "f", "default", "use header format file header_{format}")
	makeCmd.Flags().StringP("type", "t", "pdf", "save output as {type}")
	makeCmd.Flags().BoolP("landscape", "l", false, "use landscape paper orientation")
	makeCmd.Flags().BoolP("keep", "k", false, "keep generated files for debugging")
	makeCmd.Flags().BoolP("post", "", false, "generate a png for posting to social media")
	makeCmd.Flags().BoolP("root", "", false, "save result in project root")
	makeCmd.Flags().BoolP("crop", "", false, "crop page to minimal size")
	makeCmd.Flags().BoolP("point-and-click", "", false, "turn on point-and-click")
	makeCmd.Flags().BoolP("view-spacing", "", false, "turn on Paper.annotatespacing")

	lyMaker = &maker{makeCmd}
}

const fileHeader = `%% Generated from {{.sourceFile}} by domusic

\version "{{.version}}"

\pointAndClick{{if .pointAndClick}}On{{else}}Off{{end}}

#(set-global-staff-size {{.staffSize}})
#(set-default-paper-size "{{.paperSize}}" '{{if .landscape}}landscape{{else}}portrait{{end}})

\include "./bagpipe_new.ly"
\include "./bagpipe_extra.ly"
\include "./header_{{.headerFormat}}.ly"

%% Local tweaks
\paper {
  annotate-spacing = {{if .viewSpacing}}##t{{else}}##f{{end}}
  ragged-bottom = ##t
  {{if .removeTagline}}tagline = ""{{end}}
}
\layout {
  \context {
    \Score
    \override NonMusicalPaperColumn #'line-break-permission = ##f
  }
}

%% The tune to generate.
`

type maker struct {
	cmd *cobra.Command
}

func (m *maker) run(src string) error {
	var err error
	fmt.Println("Processing file", src)
	if m.flagString("type") == "pdf" {
		fmt.Println("  * Creating preview file")
		err = lyMaker.preview(src)
		if err == nil {
			fmt.Println("  * Creating PDF file")
			err = lyMaker.pdf(src)
		}
	} else {
		fmt.Println("  * Creating PNG file")
		err = lyMaker.png(src)
	}

	templateFile := getTemplatePath(src)

	if err != nil {
		e, ea, _ := getEditor()
		logFile := strings.TrimSuffix(templateFile, ".ly") + ".log"
		c := exec.Command(e, append(ea, logFile)...)
		c.Run()
		return err
	}

	if m.flagBool("keep") {
		return nil
	}

	fmt.Println("  * Cleaning up")
	lyMaker.crop(templateFile)
	cleanup(templateFile)
	if !m.flagBool("root") {
		moveFiles(templateFile, src)
	}

	return nil
}

func (m *maker) preview(src string) error {
	lyArgs := []string{
		"--png",
		"-dpreview",
		"-dno-print-pages",
		"-dresolution=84",
		"-dpreview-include-book-title",
		"-dwithout-comment",
	}

	return m.runLilypond(src, lyArgs, true)
}

func (m *maker) pdf(src string) error {
	lyArgs := []string{
		"--pdf",
	}

	return m.runLilypond(src, lyArgs, false)
}

func (m *maker) png(src string) error {
	lyArgs := []string{
		"--png",
		"-dresolution=144",
	}

	return m.runLilypond(src, lyArgs, false)
}

func (m *maker) runLilypond(src string, args []string, minimal bool) error {
	var err error

	if src != "" {
		tp, err := m.makeTemplateFile(src, minimal)
		if err != nil {
			return err
		}
		tpBase := strings.TrimSuffix(tp, ".ly")
		args = append(args, "-o"+tpBase, tp)

		c := exec.Command("lilypond", args...)
		errOut, err := c.CombinedOutput()
		ioutil.WriteFile(tpBase+".log", errOut, 0644)
	} else {
		c := exec.Command("lilypond", args...)
		err = c.Run()
	}

	return err
}

func (m *maker) crop(src string) error {
	if !m.flagBool("crop") {
		return nil
	}

	path := strings.TrimSuffix(src, ".ly") + "." + m.flagString("type")
	c := exec.Command("mogrify", "-trim", "-bordercolor", "white", "-border", "12", path)

	return c.Run()
}

func (m *maker) makeTemplateFile(sourceFile string, minimal bool) (string, error) {
	format := m.flagString("format")
	if format == "default" && strings.Contains(sourceFile, ".book") {
		format = "book"
	}
	data := map[string]interface{}{
		"sourceFile":    sourceFile,
		"version":       "2.18.0",
		"pointAndClick": m.flagBool("point-and-click"),
		"staffSize":     m.flagInt("staff-size"),
		"paperSize":     m.flagString("paper-size"),
		"landscape":     m.flagBool("landscape"),
		"headerFormat":  format,
		"viewSpacing":   m.flagBool("view-spacing"),
		"removeTagline": m.flagBool("crop"),
	}

	header, err := executeTemplate(fileHeader, data)
	if err != nil {
		errExit(err)
	}
	source, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		errExit(err)
	}

	sourceFile = ensureSuffix(noExt(sourceFile), ".ly")
	templatePath := getTemplatePath(sourceFile)
	f, err := os.OpenFile(templatePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer f.Close()

	n, err := f.WriteString(header)
	if err == nil && n < len(header) {
		return "", io.ErrShortWrite
	}

	includeLine := true
	for _, line := range bytes.Split(source, []byte("\n")) {
		trimmedLine := bytes.TrimLeft(line, " \t")
		if minimal && bytes.HasPrefix(trimmedLine, []byte("%%% START SKIP")) {
			includeLine = false
		}
		if includeLine {
			n, err = f.WriteString(string(line) + "\n")
			if err == nil && n < len(line) {
				return "", io.ErrShortWrite
			}
		}
		if minimal && bytes.HasPrefix(trimmedLine, []byte("%%% END SKIP")) {
			includeLine = true
		}
	}

	return templatePath, nil
}

func (m *maker) flagBool(name string) bool {
	val, _ := m.cmd.Flags().GetBool(name)
	return val
}

func (m *maker) flagInt(name string) int {
	val, _ := m.cmd.Flags().GetInt(name)
	return val
}

func (m *maker) flagString(name string) string {
	val, _ := m.cmd.Flags().GetString(name)
	return val
}

func cleanup(path string) {
	base := strings.TrimSuffix(path, ".ly")
	os.Remove(base + ".log")
	os.Remove(base + ".ly")
	os.Remove(base + ".preview.eps")
	os.Remove(base + ".preview.pdf")
	os.Remove(base + ".ps")
}

func moveFiles(from, to string) {
	fromBase := strings.TrimSuffix(from, ".ly")
	os.Rename(fromBase+".pdf", getPdfPath(to))
	os.Rename(fromBase+".preview.png", getPreviewPath(to))
}
