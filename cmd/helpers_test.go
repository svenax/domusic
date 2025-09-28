package cmd

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/spf13/viper"
)

func Test_getSourcePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple_path", args{"song"}, "song.ly"},
		{"path_with_extension", args{"song.ly"}, "song.ly"},
		{"path_with_subdirectory", args{"folk/song"}, "folk/song.ly"},
		{"path_with_multiple_extensions", args{"song.mid.ly"}, "song.mid.ly"},
		{"empty_path", args{""}, ".ly"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSourcePath(tt.args.path); got != tt.want {
				t.Errorf("getSourcePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getTemplatePath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple_path", args{"song.ly"}, "__song.ly"},
		{"path_with_directory", args{"/music/folk/song.ly"}, "__song.ly"},
		{"path_with_multiple_extensions", args{"song.mid.ly"}, "__song.ly"},
		{"path_without_extension", args{"song"}, "__song.ly"},
		{"empty_path", args{""}, "__.ly"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getTemplatePath(tt.args.path); got != tt.want {
				t.Errorf("getTemplatePath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPdfPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple_path", args{"song.ly"}, "_output/song.pdf"},
		{"path_with_directory", args{"folk/song.ly"}, "_output/folk/song.pdf"},
		{"path_with_multiple_extensions", args{"song.mid.ly"}, "_output/song.pdf"},
		{"path_without_extension", args{"song"}, "_output/song.pdf"},
		{"empty_path", args{""}, "_output/.pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPdfPath(tt.args.path); got != tt.want {
				t.Errorf("getPdfPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getPreviewPath(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"simple_path", args{"song.ly"}, "_output/song.preview.png"},
		{"path_with_directory", args{"folk/song.ly"}, "_output/folk/song.preview.png"},
		{"path_with_multiple_extensions", args{"song.mid.ly"}, "_output/song.preview.png"},
		{"path_without_extension", args{"song"}, "_output/song.preview.png"},
		{"empty_path", args{""}, "_output.preview.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getPreviewPath(tt.args.path); got != tt.want {
				t.Errorf("getPreviewPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOutputPath(t *testing.T) {
	type args struct {
		path    string
		preview bool
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"pdf_output", args{"song.ly", false}, "_output/song.pdf"},
		{"preview_output", args{"song.ly", true}, "_output/song.preview.png"},
		{"pdf_with_directory", args{"folk/song.ly", false}, "_output/folk/song.pdf"},
		{"preview_with_directory", args{"folk/song.ly", true}, "_output/folk/song.preview.png"},
		{"pdf_multiple_extensions", args{"song.mid.ly", false}, "_output/song.pdf"},
		{"preview_multiple_extensions", args{"song.mid.ly", true}, "_output/song.preview.png"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getOutputPath(tt.args.path, tt.args.preview); got != tt.want {
				t.Errorf("getOutputPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getEditor(t *testing.T) {
	tests := []struct {
		name     string
		lyEditor string
		editor   string
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{"ly_editor_set", "vim -n", "", "vim", []string{"-n"}, false},
		{"ly_editor_with_multiple_args", "code --wait --new-window", "", "code", []string{"--wait", "--new-window"}, false},
		{"fallback_to_editor", "", "nano", "nano", []string{}, false},
		{"fallback_to_editor_with_args", "", "emacs -nw", "emacs", []string{"-nw"}, false},
		{"no_editor_set", "", "", "", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper configuration
			viper.Reset()
			if tt.lyEditor != "" {
				viper.Set("ly-editor", tt.lyEditor)
			}
			if tt.editor != "" {
				viper.Set("editor", tt.editor)
			}

			name, args, err := getEditor()
			if (err != nil) != tt.wantErr {
				t.Errorf("getEditor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if name != tt.wantName {
				t.Errorf("getEditor() name = %v, wantName %v", name, tt.wantName)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("getEditor() args = %v, wantArgs %v", args, tt.wantArgs)
			}
		})
	}
}

func Test_getViewer(t *testing.T) {
	tests := []struct {
		name     string
		lyViewer string
		wantName string
		wantArgs []string
		wantErr  bool
	}{
		{"ly_viewer_set", "evince", "evince", []string{}, false},
		{"ly_viewer_with_args", "okular --presentation", "okular", []string{"--presentation"}, false},
		{"no_viewer_set", "", "", []string{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper configuration
			viper.Reset()
			if tt.lyViewer != "" {
				viper.Set("ly-viewer", tt.lyViewer)
			}

			name, args, err := getViewer()
			if (err != nil) != tt.wantErr {
				t.Errorf("getViewer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if name != tt.wantName {
				t.Errorf("getViewer() got = %v, wantName %v", name, tt.wantName)
			}
			if !reflect.DeepEqual(args, tt.wantArgs) {
				t.Errorf("getViewer() args = %v, wantArgs %v", args, tt.wantArgs)
			}
		})
	}
}

func Test_getNotebook(t *testing.T) {
	tests := []struct {
		name       string
		enNotebook string
		want       string
		wantErr    bool
	}{
		{"notebook_set", "Music Scores", "Music Scores", false},
		{"notebook_with_spaces", "My Music Collection", "My Music Collection", false},
		{"no_notebook_set", "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper configuration
			viper.Reset()
			if tt.enNotebook != "" {
				viper.Set("en-notebook", tt.enNotebook)
			}

			got, err := getNotebook()
			if (err != nil) != tt.wantErr {
				t.Errorf("getNotebook() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getNotebook() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_executeTemplate(t *testing.T) {
	type args struct {
		tmplString string
		data       interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			"simple_template",
			args{"Hello {{.Name}}", map[string]string{"Name": "World"}},
			"Hello World",
			false,
		},
		{
			"template_with_multiple_fields",
			args{"{{.Title}} by {{.Composer}}", map[string]string{"Title": "Symphony No. 9", "Composer": "Beethoven"}},
			"Symphony No. 9 by Beethoven",
			false,
		},
		{
			"template_with_loops",
			args{"Notes: {{range .Notes}}{{.}} {{end}}", map[string][]string{"Notes": {"C", "D", "E"}}},
			"Notes: C D E ",
			false,
		},
		{
			"empty_template",
			args{"", map[string]string{}},
			"",
			false,
		},
		{
			"template_with_no_placeholders",
			args{"Plain text", map[string]string{}},
			"Plain text",
			false,
		},
		{
			"invalid_template_syntax",
			args{"{{.Invalid.}}", map[string]string{}},
			"",
			true,
		},
		{
			"template_with_missing_field",
			args{"Hello {{.Name}}", map[string]string{}},
			"Hello <no value>",
			false,
		},
		{
			"template_with_conditional",
			args{"{{if .Show}}Visible{{else}}Hidden{{end}}", map[string]bool{"Show": true}},
			"Visible",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := executeTemplate(tt.args.tmplString, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("executeTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("executeTemplate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ensureSuffix(t *testing.T) {
	type args struct {
		path   string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"PlainNoSuffix", args{"a", ""}, "a"},
		{"PlainSuffix", args{"a", ".x"}, "a.x"},
		{"PlainHasSuffix", args{"a.x", "x"}, "a.x"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ensureSuffix(tt.args.path, tt.args.suffix); got != tt.want {
				t.Errorf("ensureSuffix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_noExt(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"PlainNoExt", args{"a"}, "a"},
		{"PathNoExt", args{"/b/a"}, "/b/a"},
		{"PlainOneExt", args{"a.b"}, "a"},
		{"PlainManyExt", args{"a.b.c.d"}, "a"},
		{"PathManyExt", args{"/e/f/a.b.c.d"}, "/e/f/a"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := noExt(tt.args.path); got != tt.want {
				t.Errorf("noExt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_makeRel(t *testing.T) {
	type args struct {
		path string
		root string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"path_within_root", args{"/music/folk/song.ly", "/music"}, "folk/song.ly"},
		{"path_within_root_with_trailing_slash", args{"/music/folk/song.ly", "/music/"}, "folk/song.ly"},
		{"path_outside_root", args{"/other/song.ly", "/music"}, "/other/song.ly"},
		{"exact_root_path", args{"/music", "/music"}, "/music"},
		{"root_with_file", args{"/music/song.ly", "/music"}, "song.ly"},
		{"empty_path", args{"", "/music"}, ""},
		{"relative_path", args{"folk/song.ly", "/music"}, "folk/song.ly"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper configuration
			viper.Reset()
			viper.Set("root", tt.args.root)

			if got := makeRel(tt.args.path); got != tt.want {
				t.Errorf("makeRel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathFromRoot(t *testing.T) {
	type args struct {
		parts []string
		root  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"single_part", args{[]string{"song.ly"}, "/music"}, "/music/song.ly"},
		{"multiple_parts", args{[]string{"folk", "song.ly"}, "/music"}, "/music/folk/song.ly"},
		{"absolute_path_unchanged", args{[]string{"/other/song.ly"}, "/music"}, "/other/song.ly"},
		{"empty_parts", args{[]string{}, "/music"}, "/music"},
		{"single_empty_part", args{[]string{""}, "/music"}, "/music"},
		{"output_directory", args{[]string{"_output", "song.pdf"}, "/music"}, "/music/_output/song.pdf"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup viper configuration
			viper.Reset()
			viper.Set("root", tt.args.root)

			if got := pathFromRoot(tt.args.parts...); got != tt.want {
				t.Errorf("pathFromRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyFile(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string // filename -> content
		srcFile       string
		dstFile       string
		wantBytes     int64
		wantErr       bool
		errorContains string
	}{
		{
			name:       "successful_copy",
			setupFiles: map[string]string{"source.txt": "Hello, World!"},
			srcFile:    "source.txt",
			dstFile:    "destination.txt",
			wantBytes:  13,
			wantErr:    false,
		},
		{
			name:       "copy_empty_file",
			setupFiles: map[string]string{"empty.txt": ""},
			srcFile:    "empty.txt",
			dstFile:    "empty_copy.txt",
			wantBytes:  0,
			wantErr:    false,
		},
		{
			name:       "copy_large_content",
			setupFiles: map[string]string{"large.txt": "This is a larger file with more content for testing purposes."},
			srcFile:    "large.txt",
			dstFile:    "large_copy.txt",
			wantBytes:  61,
			wantErr:    false,
		},
		{
			name:          "source_file_does_not_exist",
			setupFiles:    map[string]string{},
			srcFile:       "nonexistent.txt",
			dstFile:       "destination.txt",
			wantBytes:     0,
			wantErr:       true,
			errorContains: "no such file or directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tempDir, err := os.MkdirTemp("", "copyfile_test")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tempDir)

			// Setup test files
			for filename, content := range tt.setupFiles {
				fullPath := filepath.Join(tempDir, filename)
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to create test file %s: %v", filename, err)
				}
			}

			// Construct full paths
			srcPath := filepath.Join(tempDir, tt.srcFile)
			dstPath := filepath.Join(tempDir, tt.dstFile)

			// Execute the function
			got, err := copyFile(srcPath, dstPath)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("copyFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errorContains != "" {
				if err == nil {
					t.Errorf("copyFile() error = nil, want error containing %v", tt.errorContains)
					return
				}
				// Use errors.Is for os.ErrNotExist, otherwise fallback to substring
				if tt.errorContains == "no such file or directory" {
					if !errors.Is(err, os.ErrNotExist) {
						t.Errorf("copyFile() error = %v, want os.ErrNotExist", err)
					}
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("copyFile() error = %v, should contain %v", err, tt.errorContains)
				}
				return
			}

			// Check bytes copied
			if got != tt.wantBytes {
				t.Errorf("copyFile() = %v, want %v", got, tt.wantBytes)
			}

			// Verify destination file exists and has correct content (for successful copies)
			if !tt.wantErr && tt.wantBytes > 0 {
				if _, err := os.Stat(dstPath); os.IsNotExist(err) {
					t.Errorf("Destination file was not created")
				} else {
					// Verify content matches
					srcContent, _ := os.ReadFile(srcPath)
					dstContent, _ := os.ReadFile(dstPath)
					if string(srcContent) != string(dstContent) {
						t.Errorf("Destination file content doesn't match source")
					}
				}
			}
		})
	}
}

// Use strings.Contains directly in the test code instead of a helper function.
