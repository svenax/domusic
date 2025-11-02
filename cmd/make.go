package cmd

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"
)

var makeCmd = &cli.Command{
	Name:  "make",
	Usage: "Run Lilypond on music file(s)",
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:    "resolution",
			Aliases: []string{"r"},
			Value:   144,
			Usage:   "resolution for PNG files",
		},
		&cli.IntFlag{
			Name:    "staff-size",
			Aliases: []string{"s"},
			Value:   15,
			Usage:   "staff size",
		},
		&cli.StringFlag{
			Name:    "paper-size",
			Aliases: []string{"p"},
			Value:   "a4",
			Usage:   "paper size",
		},
		&cli.StringFlag{
			Name:    "format",
			Aliases: []string{"f"},
			Value:   "default",
			Usage:   "use header format file header_{format}",
		},
		&cli.StringFlag{
			Name:    "type",
			Aliases: []string{"t"},
			Value:   "pdf",
			Usage:   "save output as {type}",
		},
		&cli.BoolFlag{
			Name:    "landscape",
			Aliases: []string{"l"},
			Usage:   "use landscape paper orientation",
		},
		&cli.BoolFlag{
			Name:    "keep",
			Aliases: []string{"k"},
			Usage:   "keep generated files for debugging",
		},
		&cli.BoolFlag{
			Name:  "post",
			Usage: "generate a png for posting to social media",
		},
		&cli.BoolFlag{
			Name:  "root",
			Usage: "save result in project root",
		},
		&cli.BoolFlag{
			Name:  "crop",
			Usage: "crop page to minimal size",
		},
		&cli.BoolFlag{
			Name:  "point-and-click",
			Usage: "turn on point-and-click",
		},
		&cli.BoolFlag{
			Name:  "view-spacing",
			Usage: "turn on Paper.annotatespacing",
		},
		&cli.StringFlag{
			Name:  "font-include",
			Usage: "include font configuration file",
		},
	},

	Action: func(ctx context.Context, cmd *cli.Command) error {
		// Set font-include in config
		setString("font-include", cmd.String("font-include"))

		args := cmd.Args().Slice()
		maker := &maker{cmd}
		for _, arg := range args {
			files := []string{arg}
			if strings.Contains(arg, "*") {
				files, _ = filepath.Glob(pathFromRoot(arg))
			}
			for _, f := range files {
				maker.run(getSourcePath(f))
			}
		}
		return nil
	},
}

const fileHeader = `%% Generated from {{.sourceFile}} by domusic

\version "{{.version}}"

\pointAndClick{{if .pointAndClick}}On{{else}}Off{{end}}

#(set-global-staff-size {{.staffSize}})
#(set-default-paper-size "{{.paperSize}}" '{{if .landscape}}landscape{{else}}portrait{{end}})

{{if .fontInclude}}\include "{{.fontInclude}}.ily"{{end}}
\include "bagpipe.ly"
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
    \override NonMusicalPaperColumn.line-break-permission = ##f
  }
}

%% The tune to generate.
`

type maker struct {
	cmd *cli.Command
}

func (m *maker) run(src string) error {
	var err error
	fmt.Println("Processing file", src)

	// Handle post flag overrides
	outputType := m.cmd.String("type")
	resolution := m.cmd.Int("resolution")
	if m.cmd.Bool("post") {
		outputType = "png"
		if resolution == 144 { // default resolution
			resolution = 84
		}
	}

	if outputType == "pdf" {
		fmt.Println("  * Creating preview file")
		err = m.preview(src, resolution)
		if err == nil {
			fmt.Println("  * Creating PDF file")
			err = m.pdf(src)
		}
	} else {
		fmt.Println("  * Creating PNG file")
		err = m.png(src, resolution)
	}

	templateFile := getTemplatePath(src)

	if err != nil {
		fmt.Println("  * Opening log file")
		e, ea, _ := getEditor()
		logFile := strings.TrimSuffix(templateFile, ".ly") + ".log"
		c := exec.Command(e, append(ea, logFile)...)
		c.Run()
		return err
	}

	if m.cmd.Bool("keep") {
		return nil
	}

	fmt.Println("  * Cleaning up")
	if m.cmd.Bool("crop") || m.cmd.Bool("post") {
		m.crop(templateFile)
	}
	cleanup(templateFile)
	if !m.cmd.Bool("root") {
		moveFiles(templateFile, src)
	}

	return nil
}

func (m *maker) preview(src string, resolution int) error {
	lyArgs := []string{
		"--png",
		"-dpreview",
		"-dno-print-pages",
		fmt.Sprintf("-dresolution=%d", resolution),
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

func (m *maker) png(src string, resolution int) error {
	lyArgs := []string{
		"--png",
		fmt.Sprintf("-dresolution=%d", resolution),
	}

	return m.runLilypond(src, lyArgs, false)
}

func (m *maker) runLilypond(src string, args []string, minimal bool) error {
	if src != "" {
		tp, err := m.makeTemplateFile(src, minimal)
		if err != nil {
			return err
		}
		tpBase := strings.TrimSuffix(tp, ".ly")
		args = append(args, "-o"+tpBase, tp)

		c := exec.Command("lilypond", args...)
		errOut, err := c.CombinedOutput()
		os.WriteFile(tpBase+".log", errOut, 0644)
		return err
	}
	c := exec.Command("lilypond", args...)

	return c.Run()
}

func (m *maker) crop(src string) error {
	path := strings.TrimSuffix(src, ".ly") + ".png"
	c := exec.Command("mogrify", "-trim", "-bordercolor", "white", "-border", "12", path)

	return c.Run()
}

func (m *maker) makeTemplateFile(sourceFile string, minimal bool) (string, error) {
	format := m.cmd.String("format")
	if format == "default" && strings.Contains(sourceFile, ".book") {
		format = "book"
	}
	data := map[string]any{
		"sourceFile":    sourceFile,
		"version":       "2.24.0",
		"pointAndClick": m.cmd.Bool("point-and-click"),
		"staffSize":     m.cmd.Int("staff-size"),
		"paperSize":     m.cmd.String("paper-size"),
		"landscape":     m.cmd.Bool("landscape"),
		"headerFormat":  format,
		"viewSpacing":   m.cmd.Bool("view-spacing"),
		"removeTagline": m.cmd.Bool("crop") || m.cmd.Bool("post"),
		"fontInclude":   getString("font-include"),
	}

	header, err := executeTemplate(fileHeader, data)
	errExit(err)
	source, err := os.ReadFile(sourceFile)
	errExit(err)

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

func cleanup(path string) {
	base := strings.TrimSuffix(path, ".ly")
	// Ignore errors for cleanup operations as files may not exist
	_ = os.Remove(base + ".log")
	_ = os.Remove(base + ".ly")
	_ = os.Remove(base + ".preview.eps")
	_ = os.Remove(base + ".preview.pdf")
	_ = os.Remove(base + ".ps")
}

func moveFiles(from, to string) {
	fromBase := strings.TrimSuffix(from, ".ly")
	os.MkdirAll(filepath.Dir(getPdfPath(to)), 0755) // Both files go in the same directory
	os.Rename(fromBase+".pdf", getPdfPath(to))
	os.Rename(fromBase+".preview.png", getPreviewPath(to))
}
