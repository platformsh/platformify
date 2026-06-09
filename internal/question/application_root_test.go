package question

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/internal/question/models"
)

func TestApplicationRoot_Ask(t *testing.T) {
	tests := []struct {
		name               string
		fileSystem         fs.FS
		dependencyManagers []models.DepManager
		want               string
		wantErr            bool
	}{
		{
			name: "Requirements.txt root",
			fileSystem: fstest.MapFS{
				"requirements.txt": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Pip},
			want:               ".",
			wantErr:            false,
		},
		{
			name: "Requirements.txt subdir",
			fileSystem: fstest.MapFS{
				"sub/requirements.txt": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Pip},
			want:               "sub",
			wantErr:            false,
		},
		{
			name: "Package.json root",
			fileSystem: fstest.MapFS{
				"package.json": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Yarn},
			want:               ".",
			wantErr:            false,
		},
		{
			name: "Package.json subdir",
			fileSystem: fstest.MapFS{
				"sub/package.json": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Npm},
			want:               "sub",
			wantErr:            false,
		},
		{
			name: "pyproject.toml subdir",
			fileSystem: fstest.MapFS{
				"sub/pyproject.toml": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Poetry},
			want:               "sub",
			wantErr:            false,
		},
		{
			name: "Pipfile subdir",
			fileSystem: fstest.MapFS{
				"sub/Pipfile": &fstest.MapFile{},
			},
			dependencyManagers: []models.DepManager{models.Pipenv},
			want:               "sub",
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &ApplicationRoot{}
			a := models.NewAnswers()
			a.WorkingDirectory = tt.fileSystem
			a.DependencyManagers = tt.dependencyManagers
			ctx := models.ToContext(context.Background(), a)

			if err := q.Ask(ctx); (err != nil) != tt.wantErr {
				t.Errorf("ApplicationRoot.Ask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && a.ApplicationRoot != tt.want {
				t.Errorf("ApplicationRoot.Ask().ApplicationRoot = %s, want %s", a.ApplicationRoot, tt.want)
			}
		})
	}
}
