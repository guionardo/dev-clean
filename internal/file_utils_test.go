package internal

import (
	"os"
	"path"
	"reflect"
	"sort"
	"testing"
)

func CreateFile(t *testing.T, filename string, dir ...string) (err error) {
	d := path.Join(dir...)
	if err := os.Mkdir(d, 0777); err != nil && !os.IsExist(err) {
		t.Error(err)
	}
	filename = path.Join(d, filename)
	f, err := os.Create(filename)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Created %s\n", filename)
	return f.Close()
}

func TestGetFiles(t *testing.T) {
	tmpDir := t.TempDir()

	type args struct {
		basePath string
		masks    []string
	}
	tests := []struct {
		name           string
		args           args
		wantSubFolders []string
		wantFileNames  []string
	}{
		{
			name: "Local",
			args: args{masks: []string{"*.go", "go.*"}},
			wantSubFolders: []string{
				path.Join(tmpDir, "dir_0/dir_01"),
				path.Join(tmpDir, "dir_0"),
				path.Join(tmpDir, "dir_1")},

			wantFileNames: []string{
				path.Join(tmpDir, "/dir_0/dir_01/test_4.go"),
				path.Join(tmpDir, "/dir_0/test_1.go"),
				path.Join(tmpDir, "/dir_1/test_5.go"),
				path.Join(tmpDir, "/go.mod"),
				path.Join(tmpDir, "/test_0.go"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CreateFile(t, "test_0.go", tmpDir)
			CreateFile(t, "test_1.go", tmpDir, "dir_0")
			CreateFile(t, "test_2.py", tmpDir, "dir_0")
			CreateFile(t, "go.mod", tmpDir)
			CreateFile(t, "test_4.go", tmpDir, "dir_0", "dir_01")
			CreateFile(t, "test_5.go", tmpDir, "dir_1")
			gotSubFolders, gotFileNames, _, err := GetFiles(tmpDir, tt.args.masks...)
			if err != nil {
				t.Error(err)
			}
			sort.Strings(gotSubFolders)
			sort.Strings(gotFileNames)
			sort.Strings(tt.wantSubFolders)
			sort.Strings(tt.wantFileNames)

			if !reflect.DeepEqual(gotSubFolders, tt.wantSubFolders) {
				t.Errorf("GetFiles() gotSubFolders = %v, want %v", gotSubFolders, tt.wantSubFolders)
			}
			if !reflect.DeepEqual(gotFileNames, tt.wantFileNames) {
				t.Errorf("GetFiles() gotFileNames = %v, want %v", gotFileNames, tt.wantFileNames)
			}
		})
	}
}
