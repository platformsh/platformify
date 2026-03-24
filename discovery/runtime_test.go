package discovery

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/platformsh/platformify/platformifier"
)

func TestDiscoverer_discoverType(t *testing.T) {
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
			name: "Python",
			fields: fields{
				fileSystem: fstest.MapFS{
					"hey.py":     &fstest.MapFile{},
					"hey.js":     &fstest.MapFile{},
					"another.py": &fstest.MapFile{},
				},
			},
			want:    "python",
			wantErr: false,
		},
		{
			name: "Express stack override",
			fields: fields{
				fileSystem: fstest.MapFS{
					"hey.py":     &fstest.MapFile{},
					"hey.js":     &fstest.MapFile{},
					"another.py": &fstest.MapFile{},
				},
				memory: map[string]any{"stack": platformifier.Express},
			},
			want:    "nodejs",
			wantErr: false,
		},
		{
			name: "Node.js skip vendor",
			fields: fields{
				fileSystem: fstest.MapFS{
					"vendor/a.py": &fstest.MapFile{},
					"vendor/b.py": &fstest.MapFile{},
					"vendor/c.py": &fstest.MapFile{},
					"hey.js":      &fstest.MapFile{},
				},
			},
			want:    "nodejs",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.fields.memory == nil {
				tt.fields.memory = make(map[string]any)
			}
			d := &Discoverer{
				fileSystem: tt.fields.fileSystem,
				memory:     tt.fields.memory,
			}
			got, err := d.discoverType()
			if (err != nil) != tt.wantErr {
				t.Errorf("Discoverer.discoverType() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Discoverer.discoverType() = %v, want %v", got, tt.want)
			}
		})
	}
}
