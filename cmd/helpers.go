package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	pth "path"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// errExit prints an error message and then exits.
func errExit(msg interface{}) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}

// getSourcePath returns the full path to a Lilypond file in the music
// hierarchy. It does not check if the file exists.
func getSourcePath(path string) string {
	return ensureSuffix(pathFromRoot(makeRel(path)), ".ly")
}

// getTemplatePath returns the full path to a generated template file in
// the root directory of the music hierarchy.
func getTemplatePath(path string) string {
	return ensureSuffix(pathFromRoot("__"+noExt(pth.Base(path))), ".ly")
}

// getPdfPath returns the full path to where the PDF file for a given tune
// should be stored.
func getPdfPath(path string) string {
	return ensureSuffix(pathFromRoot("_output", noExt(makeRel(path))), ".pdf")
}

// getPreviewPath returns the full path to where the preview file for a given
// tune should be stored.
func getPreviewPath(path string) string {
	return ensureSuffix(pathFromRoot("_output", noExt(makeRel(path))), ".preview.png")
}

// getOutputPath returns the full path to either preview or PDF file depending
// on the flag {preview}.
func getOutputPath(path string, preview bool) string {
	if preview {
		return getPreviewPath(path)
	}
	return getPdfPath(path)
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

// getViewer returns the editor set in the configuration file or exported from
// the shell. It returns the editor name and an array of arguments so it can
// easily be slotted in to exec.Command.
func getViewer() (string, []string, error) {
	var cmds []string
	if viper.IsSet("ly-editor") {
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
func executeTemplate(tmplString string, data interface{}) (string, error) {
	tmpl, err := template.New("").Parse(tmplString)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	return buf.String(), err
}

// ensureSuffix ensures the given path ends with {suffix}.
func ensureSuffix(path string, suffix string) string {
	return strings.TrimSuffix(path, suffix) + suffix
}

// noExt strips all extensions from the given path.
// E.g. given '/a/b.c.d' it will return '/a/b'.
func noExt(path string) string {
	dir := pth.Dir(path)
	base := pth.Base(path)
	if i := strings.Index(base, "."); i > 0 {
		base = base[:i]
	}

	return pth.Join(dir, base)
}

// makeRel returns a relative path from the music root. If the path is not
// actually within the music root, it will be returned as is.
func makeRel(path string) string {
	return strings.TrimPrefix(path, ensureSuffix(viper.GetString("root"), "/"))
}

// pathFromRoot returns an absolute path starting with music root and then
// all the given parts. If the given path is already absolute, it is returned
// as is.
func pathFromRoot(parts ...string) string {
	fullPath := pth.Join(parts...)
	if !pth.IsAbs(fullPath) {
		fullPath = pth.Join(viper.GetString("root"), fullPath)
	}

	return fullPath
}
