package discovery

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

func TestDiscoverer_discoverEnvironment(t *testing.T) {
	type fields struct {
		fileSystem fs.FS
		memory     map[string]any
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr bool
	}{
		{
			name: "No env",
			fields: fields{
				fileSystem: fstest.MapFS{},
				memory:     map[string]any{},
			},
			want: make(map[string]string),
		},
		{
			name: "Simple poetry",
			fields: fields{
				fileSystem: fstest.MapFS{},
				memory: map[string]any{
					"dependency_managers": []string{"poetry"},
				},
			},
			want: map[string]string{
				"POETRY_VERSION":                "1.4.0",
				"POETRY_VIRTUALENVS_IN_PROJECT": "true",
			},
		},
		{
			name: "Node.js",
			fields: fields{
				fileSystem: fstest.MapFS{},
				memory: map[string]any{
					"dependency_managers": []string{"yarn"},
					"type":                "nodejs",
				},
			},
			want: map[string]string{},
		},
		{
			name: "Node.js on different runtime",
			fields: fields{
				fileSystem: fstest.MapFS{},
				memory: map[string]any{
					"dependency_managers": []string{"yarn"},
					"type":                "python",
				},
			},
			want: map[string]string{
				"N_PREFIX": "/app/.global",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Discoverer{
				fileSystem: tt.fields.fileSystem,
				memory:     tt.fields.memory,
			}
			got, err := d.discoverEnvironment()
			if (err != nil) != tt.wantErr {
				t.Errorf("Discoverer.discoverEnvironment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Discoverer.discoverEnvironment() = %v, want %v", got, tt.want)
			}
		})
	}
}
