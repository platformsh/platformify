package discovery

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/platformifier"
)

func TestDiscoverer_discoverBuildSteps(t *testing.T) {
	type fields struct {
		fileSystem fs.FS
		memory     map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    []string
		wantErr bool
	}{
		{
			name: "Poetry Django",
			fields: fields{
				fileSystem: fstest.MapFS{
					"project/manage.py": &fstest.MapFile{},
				},
				memory: map[string]any{
					"stack":               platformifier.Django,
					"type":                "python",
					"dependency_managers": []string{"poetry"},
					"application_root":    ".",
				},
			},
			want: []string{
				"# Set PIP_USER to 0 so that Poetry does not complain",
				"export PIP_USER=0",
				"# Install poetry as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$POETRY_VERSION",
				"poetry install",
				"# Collect static files",
				"poetry run python project/manage.py collectstatic --noinput",
			},
		},
		{
			name: "Pipenv Django with Yarn build",
			fields: fields{
				fileSystem: fstest.MapFS{
					"project/manage.py": &fstest.MapFile{},
					"package.json":      &fstest.MapFile{Data: []byte(`{"scripts": {"build": "nuxt build"}}`)},
				},
				memory: map[string]any{
					"stack":               platformifier.Django,
					"type":                "python",
					"dependency_managers": []string{"poetry", "yarn"},
					"application_root":    ".",
				},
			},
			want: []string{
				"n auto || n lts",
				"hash -r",
				"yarn",
				"yarn build",
				"# Set PIP_USER to 0 so that Poetry does not complain",
				"export PIP_USER=0",
				"# Install poetry as a global tool",
				"python -m venv /app/.global",
				"pip install poetry==$POETRY_VERSION",
				"poetry install",
				"# Collect static files",
				"poetry run python project/manage.py collectstatic --noinput",
			},
		},
		{
			name: "Next.js without build script",
			fields: fields{
				fileSystem: fstest.MapFS{},
				memory: map[string]any{
					"stack":               platformifier.NextJS,
					"type":                "nodejs",
					"dependency_managers": []string{"npm"},
					"application_root":    ".",
				},
			},
			want: []string{
				"npm i",
				"npm exec next build",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				fileSystem: tt.fields.fileSystem,
				memory:     tt.fields.memory,
			}
			got, err := d.discoverBuildSteps()
			if (err != nil) != tt.wantErr {
				t.Errorf("Discoverer.discoverBuildSteps() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Discoverer.discoverBuildSteps() = %v, want %v", got, tt.want)
			}
		})
	}
}
