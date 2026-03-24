package discovery

import (
	"io/fs"
	"reflect"
	"slices"
	"testing"
	"testing/fstest"
)

func TestDiscoverer_discoverDependencyManagers(t *testing.T) {
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
			name: "Simple",
			fields: fields{
				fileSystem: fstest.MapFS{
					"package-lock.json": &fstest.MapFile{},
				},
			},
			want:    []string{"npm"},
			wantErr: false,
		},
		{
			name: "Multiple",
			fields: fields{
				fileSystem: fstest.MapFS{
					"package-lock.json": &fstest.MapFile{},
					"poetry.lock":       &fstest.MapFile{},
				},
			},
			want:    []string{"npm", "poetry"},
			wantErr: false,
		},
		{
			name: "Priority",
			fields: fields{
				fileSystem: fstest.MapFS{
					"package-lock.json": &fstest.MapFile{},
					"poetry.lock":       &fstest.MapFile{},
					"requirements.txt":  &fstest.MapFile{},
				},
			},
			want:    []string{"npm", "poetry"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				fileSystem: tt.fields.fileSystem,
				memory:     tt.fields.memory,
			}
			got, err := d.discoverDependencyManagers()
			slices.Sort(got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Discoverer.discoverDependencyManagers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Discoverer.discoverDependencyManagers() = %v, want %v", got, tt.want)
			}
		})
	}
}
