package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

var titleRx = regexp.MustCompile("title\\s*=\\s*\"(.+)\"")

// msgExit prints an error message and exits.
func msgExit(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

// errExit prints an error message and exits if err is not nil.
func errExit(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// getSourcePath returns the full path to a Lilypond file in the music
// hierarchy. It does not check if the file exists.
func getSourcePath(p string) string {
	return ensureSuffix(pathFromRoot(makeRel(p)), ".ly")
}

// getTemplatePath returns the full path to a generated template file in
// the root directory of the music hierarchy.
func getTemplatePath(p string) string {
	return ensureSuffix(pathFromRoot("__"+noExt(path.Base(p))), ".ly")
}

// getPdfPath returns the full path to where the PDF file for a given tune
// should be stored.
func getPdfPath(p string) string {
	return ensureSuffix(pathFromRoot("_output", noExt(makeRel(p))), ".pdf")
}

// getPreviewPath returns the full path to where the preview file for a given
// tune should be stored.
func getPreviewPath(p string) string {
	return ensureSuffix(pathFromRoot("_output", noExt(makeRel(p))), ".preview.png")
}

// getOutputPath returns the full path to either preview or PDF file depending
// on the flag {preview}.
func getOutputPath(p string, preview bool) string {
	if preview {
		return getPreviewPath(p)
	}
	return getPdfPath(p)
}

// getEditor returns the editor set in the configuration file or exported from
// the shell. It returns the editor name and an array of arguments so it can
// easily be slotted in to exec.Command.
func getEditor() (string, []string, error) {
	var cmds []string
	if viper.IsSet("ly-editor") {
		cmds = strings.Split(viper.GetString("ly-editor"), " ")
		return cmds[0], cmds[1:], nil
	}
	if viper.IsSet("editor") {
		cmds = strings.Split(viper.GetString("editor"), " ")
		return cmds[0], cmds[1:], nil
	}
	return "", []string{}, errors.New("no editor set")
}

// getViewer returns the viewer set in the configuration file or exported from
// the shell. It returns the viewer name and an array of arguments so it can
// easily be slotted in to exec.Command.
func getViewer() (string, []string, error) {
	var cmds []string
	if viper.IsSet("ly-viewer") {
		cmds = strings.Split(viper.GetString("ly-viewer"), " ")
		return cmds[0], cmds[1:], nil
	}
	return "", []string{}, errors.New("no PDF viewer set")
}

// getNotebook returns the Evernote notebook set in the configuration file or
// exported from the shell.
func getNotebook() (string, error) {
	if viper.IsSet("en-notebook") {
		return viper.GetString("en-notebook"), nil
	}
	return "", errors.New("no notebook set")
}

// executeTemplate takes a text template and a data map, and returns the
// text with the data inserted.
func executeTemplate(tmplString string, data any) (string, error) {
	tmpl, err := template.New("").Parse(tmplString)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	return buf.String(), err
}

// ensureSuffix ensures the given path ends with {suffix}.
func ensureSuffix(p string, suffix string) string {
	return strings.TrimSuffix(p, suffix) + suffix
}

// noExt strips all extensions from the given path.
// E.g. given '/a/b.c.d' it will return '/a/b'.
func noExt(p string) string {
	dir, base := path.Split(p)
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}

	return path.Join(dir, base)
}

// makeRel returns a relative path from the music root. If the path is not
// actually within the music root, it will be returned as is.
func makeRel(p string) string {
	return strings.TrimPrefix(p, ensureSuffix(viper.GetString("root"), "/"))
}

// pathFromRoot returns an absolute path starting with music root and then
// all the given parts. If the given path is already absolute, it is returned
// as is.
func pathFromRoot(parts ...string) string {
	fullPath := path.Join(parts...)
	if !path.IsAbs(fullPath) {
		fullPath = path.Join(viper.GetString("root"), fullPath)
	}

	return fullPath
}

// copyFile copies a file from src to dst. It returns the number of bytes
// copied and an error status.
func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
