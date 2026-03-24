package question

import (
	"context"
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/discovery"
	"github.com/platformsh/platformify/internal/question/models"
)

func TestStack_Ask(t *testing.T) {
	tests := []struct {
		name       string
		fileSystem fs.FS
		want       string
		wantErr    bool
	}{
		{
			name: "Django",
			fileSystem: fstest.MapFS{
				"demo/settings.py": &fstest.MapFile{},
				"manage.py":        &fstest.MapFile{},
			},
			want:    "Django",
			wantErr: false,
		},
		{
			name: "Django subdir",
			fileSystem: fstest.MapFS{
				"sub/settings.py":    &fstest.MapFile{},
				"sub/demo/manage.py": &fstest.MapFile{},
			},
			want:    "Django",
			wantErr: false,
		},
		{
			name: "Flask requirements.txt",
			fileSystem: fstest.MapFS{
				"requirements.txt": &fstest.MapFile{Data: []byte("FlAsK==1.2.3#hash-here")},
			},
			want:    "Flask",
			wantErr: false,
		},
		{
			name: "Flask Poetry",
			fileSystem: fstest.MapFS{
				"pyproject.toml": &fstest.MapFile{Data: []byte(`
[tool.poetry.dependencies]
# Get the latest revision on the branch named "next"
requests = { git = "https://github.com/kennethreitz/requests.git", branch = "next" }
# Get a revision by its commit hash
FlAsK = { git = "https://github.com/pallets/flask.git", rev = "38eb5d3b" }
				`)},
			},
			want:    "Flask",
			wantErr: false,
		},
		{
			name: "Flask Pipenv",
			fileSystem: fstest.MapFS{
				"Pipfile": &fstest.MapFile{Data: []byte(`
[packages]
fLaSk = "^1.2.3"
				`)},
			},
			want:    "Flask",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Stack{}
			a := models.NewAnswers()
			a.WorkingDirectory = tt.fileSystem
			a.Discoverer = discovery.New(tt.fileSystem)
			ctx := models.ToContext(context.Background(), a)

			if err := q.Ask(ctx); (err != nil) != tt.wantErr {
				t.Errorf("Stack.Ask() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && a.Stack.Title() != tt.want {
				t.Errorf("Stack.Ask().Stack = %s, want %s", a.Stack.Title(), tt.want)
			}
		})
	}
}
