package discovery

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestDiscoverer_discoverApplicationRoot(t *testing.T) {
	type fields struct {
		fileSystem fs.FS
		memory     map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Simple",
			fields: fields{
				fileSystem: fstest.MapFS{
					"package-lock.json": &fstest.MapFile{},
				},
			},
			want: ".",
		},
		{
			name: "No root",
			fields: fields{
				fileSystem: fstest.MapFS{},
			},
			want: "",
		},
		{
			name: "Priority",
			fields: fields{
				fileSystem: fstest.MapFS{
					"yarn/yarn.lock":              &fstest.MapFile{},
					"poetry/poetry.lock":          &fstest.MapFile{},
					"composer/composer-lock.json": &fstest.MapFile{},
				},
			},
			want: "poetry",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				fileSystem: tt.fields.fileSystem,
				memory:     tt.fields.memory,
			}
			if d.memory == nil {
				d.memory = make(map[string]any)
			}
			got, err := d.discoverApplicationRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("Discoverer.discoverApplicationRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Discoverer.discoverApplicationRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
