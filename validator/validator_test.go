package validator

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func Test_validateUpsunConfig(t *testing.T) {
	type args struct {
		path fs.FS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1:
    type: "python:3.11"
`),
					},
				},
			},
		},
		{
			name: "marshal error",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:a
  app1:
    type: "python:3.11"
`),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "sub-directory",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1:
    type: "python:3.11"
`,
						),
					},
					".upsun/local/config.yaml": &fstest.MapFile{
						Data: []byte(`invalid:yaml`),
					},
				},
			},
			wantErr: false,
		},
		{
			name: "clashing key",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1:
    type: "python:3.11"
`,
						),
					},
					".upsun/another.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1:
    type: "python:3.11"
`,
						),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing key",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1: {}
`,
						),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple files",
			args: args{
				path: fstest.MapFS{
					".upsun/config.yaml": &fstest.MapFile{
						Data: []byte(`
applications:
  app1:
    type: "python:3.11"
`,
						),
					},
					".upsun/services.yaml": &fstest.MapFile{
						Data: []byte(`
services:
  redis:
    type: redis:6.2
    size: AUTO
`,
						),
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateUpsunConfig(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("validateUpsunConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
