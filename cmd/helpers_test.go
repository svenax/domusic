package cmd

import (
	"reflect"
	"testing"
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		name    string
		want    string
		want1   []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getEditor()
			if (err != nil) != tt.wantErr {
				t.Errorf("getEditor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getEditor() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getEditor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getViewer(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		want1   []string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := getViewer()
			if (err != nil) != tt.wantErr {
				t.Errorf("getViewer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getViewer() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("getViewer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_getNotebook(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		// TODO: Add test cases.
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
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := makeRel(tt.args.path); got != tt.want {
				t.Errorf("makeRel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pathFromRoot(t *testing.T) {
	type args struct {
		parts []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathFromRoot(tt.args.parts...); got != tt.want {
				t.Errorf("pathFromRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
